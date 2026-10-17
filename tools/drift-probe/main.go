// Command drift-probe measures a deployed server's clock offset against NTP over /ws.
package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/beevik/ntp"
	"github.com/gorilla/websocket"
)

// Fixed subscribe location (Copenhagen); the probe only reads utc.
const (
	subscribeLat = 55.6761
	subscribeLon = 12.5683
)

const (
	dialTimeout   = 60 * time.Second
	sampleTimeout = 10 * time.Second
)

var csvHeader = []string{
	"run_at", "wake_ms", "ntp_offset_ms", "median_rtt_ms",
	"median_abs_offset_ms", "max_abs_offset_ms", "samples",
}

type update struct {
	UTC   time.Time `json:"utc"`
	Error string    `json:"error"`
}

// result is one probe run.
type result struct {
	RunAt     time.Time
	Wake      time.Duration
	NTPOffset time.Duration
	Offsets   []time.Duration
	RTTs      []time.Duration
}

func main() {
	url := flag.String("url", "", "WebSocket URL, e.g. wss://<app>.fly.dev/ws (required)")
	ntpHost := flag.String("ntp", "pool.ntp.org", "NTP server")
	samples := flag.Int("samples", 10, "number of updates to measure")
	out := flag.String("out", "", "CSV file to append the run to")
	flag.Parse()

	if *url == "" || *samples < 1 {
		flag.Usage()
		os.Exit(2)
	}

	r, err := probe(*url, *ntpHost, *samples)
	if err != nil {
		log.Fatalf("drift-probe: %v", err)
	}

	fmt.Printf("wake %v, ntp offset %v, median rtt %v, median |offset| %v, max |offset| %v\n",
		r.Wake.Round(time.Millisecond), r.NTPOffset.Round(time.Microsecond),
		median(r.RTTs).Round(time.Microsecond),
		median(absAll(r.Offsets)).Round(time.Microsecond), maxAbs(r.Offsets).Round(time.Microsecond))

	if *out != "" {
		if err := appendCSV(*out, r); err != nil {
			log.Fatalf("drift-probe: %v", err)
		}
	}
}

// probe corrects the local clock with NTP, then compares server utc against the RTT midpoint of each sample.
func probe(url, ntpHost string, samples int) (result, error) {
	r := result{RunAt: time.Now().UTC()}

	ntpResp, err := ntp.Query(ntpHost)
	if err != nil {
		return r, fmt.Errorf("query ntp %s: %w", ntpHost, err)
	}
	if err := ntpResp.Validate(); err != nil {
		return r, fmt.Errorf("validate ntp response: %w", err)
	}
	r.NTPOffset = ntpResp.ClockOffset

	dialer := websocket.Dialer{HandshakeTimeout: dialTimeout}
	start := time.Now()
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return r, fmt.Errorf("dial %s: %w", url, err)
	}
	r.Wake = time.Since(start)
	defer func() { _ = conn.Close() }()

	if err := conn.WriteJSON(map[string]float64{"lat": subscribeLat, "lon": subscribeLon}); err != nil {
		return r, fmt.Errorf("send subscribe: %w", err)
	}

	// The pong handler runs inside ReadMessage, on this goroutine.
	var pongAt time.Time
	conn.SetPongHandler(func(string) error {
		pongAt = time.Now()
		return nil
	})

	for range samples {
		pongAt = time.Time{}
		pingAt := time.Now()
		if err := conn.WriteControl(websocket.PingMessage, nil, pingAt.Add(sampleTimeout)); err != nil {
			return r, fmt.Errorf("send ping: %w", err)
		}
		if err := conn.SetReadDeadline(pingAt.Add(sampleTimeout)); err != nil {
			return r, fmt.Errorf("set read deadline: %w", err)
		}

		// Skip updates that arrive before the pong; use the first one after it.
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return r, fmt.Errorf("read update: %w", err)
			}
			recv := time.Now()
			if pongAt.IsZero() {
				continue
			}

			var u update
			if err := json.Unmarshal(data, &u); err != nil {
				return r, fmt.Errorf("decode update: %w", err)
			}
			if u.Error != "" {
				return r, fmt.Errorf("server rejected subscribe: %s", u.Error)
			}
			if u.UTC.IsZero() {
				return r, errors.New("update has no utc field")
			}

			rtt := pongAt.Sub(pingAt)
			r.RTTs = append(r.RTTs, rtt)
			r.Offsets = append(r.Offsets, offset(u.UTC, recv, rtt, r.NTPOffset))
			break
		}
	}

	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	_ = conn.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second))
	return r, nil
}

// offset is the server clock's error: utc minus the NTP-corrected moment the server sent the update.
func offset(utc, recv time.Time, rtt, ntpOffset time.Duration) time.Duration {
	sent := recv.Add(ntpOffset).Add(-rtt / 2)
	return utc.Sub(sent)
}

func absAll(ds []time.Duration) []time.Duration {
	out := make([]time.Duration, len(ds))
	for i, d := range ds {
		out[i] = d.Abs()
	}
	return out
}

func median(ds []time.Duration) time.Duration {
	if len(ds) == 0 {
		return 0
	}
	s := slices.Clone(ds)
	slices.Sort(s)
	n := len(s)
	if n%2 == 1 {
		return s[n/2]
	}
	return (s[n/2-1] + s[n/2]) / 2
}

func maxAbs(ds []time.Duration) time.Duration {
	var m time.Duration
	for _, d := range ds {
		m = max(m, d.Abs())
	}
	return m
}

// appendCSV appends r to path, writing the header first when the file is new or empty.
func appendCSV(path string, r result) error {
	info, err := os.Stat(path)
	needHeader := errors.Is(err, fs.ErrNotExist) || (err == nil && info.Size() == 0)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	f, err := os.OpenFile(filepath.Clean(path), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}

	w := csv.NewWriter(f)
	rows := [][]string{csvRow(r)}
	if needHeader {
		rows = append([][]string{csvHeader}, rows...)
	}
	if err := w.WriteAll(rows); err != nil {
		_ = f.Close()
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close %s: %w", path, err)
	}
	return nil
}

func csvRow(r result) []string {
	return []string{
		r.RunAt.Format(time.RFC3339),
		ms(r.Wake),
		ms(r.NTPOffset),
		ms(median(r.RTTs)),
		ms(median(absAll(r.Offsets))),
		ms(maxAbs(r.Offsets)),
		strconv.Itoa(len(r.Offsets)),
	}
}

func ms(d time.Duration) string {
	return strconv.FormatFloat(float64(d)/float64(time.Millisecond), 'f', 3, 64)
}
