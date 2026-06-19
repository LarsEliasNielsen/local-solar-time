package server

import (
	"net/http/httptest"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"local-solar-time/internal/clock"
	"local-solar-time/internal/solar"
)

func wsURL(ts *httptest.Server) string {
	return "ws" + strings.TrimPrefix(ts.URL, "http")
}

func TestServeIndependentClients(t *testing.T) {
	fixed := clock.FixedClock{Time: time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)}
	s := New(fixed, 10*time.Millisecond)
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	for i := range 10 {
		lat, lon := float64(i), float64(i)*2

		conn, _, err := websocket.DefaultDialer.Dial(wsURL(ts), nil)
		if err != nil {
			t.Fatalf("dial: %v", err)
		}
		func() {
			defer func() { _ = conn.Close() }()

			if err := conn.WriteJSON(SubscribeRequest{Lat: lat, Lon: lon}); err != nil {
				t.Fatalf("write subscribe: %v", err)
			}

			var got UpdateResponse
			if err := conn.ReadJSON(&got); err != nil {
				t.Fatalf("read update: %v", err)
			}

			want := toUpdateResponse(solar.Compute(fixed.Now(), lat, lon))
			if !reflect.DeepEqual(got, want) {
				t.Errorf("lat=%v lon=%v: got %+v, want %+v", lat, lon, got, want)
			}
		}()
	}
}

func TestServeRejectsInvalidCoordinatesAndAllowsRetry(t *testing.T) {
	fixed := clock.FixedClock{Time: time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)}
	s := New(fixed, 10*time.Millisecond)
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	conn, _, err := websocket.DefaultDialer.Dial(wsURL(ts), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	if err := conn.WriteJSON(SubscribeRequest{Lat: 200, Lon: 0}); err != nil {
		t.Fatalf("write invalid subscribe: %v", err)
	}

	var errResp ErrorResponse
	if err := conn.ReadJSON(&errResp); err != nil {
		t.Fatalf("read error response: %v", err)
	}
	if errResp.Error == "" {
		t.Error("expected a non-empty error message")
	}

	if err := conn.WriteJSON(SubscribeRequest{Lat: 10, Lon: 10}); err != nil {
		t.Fatalf("write valid subscribe after retry: %v", err)
	}

	var got UpdateResponse
	if err := conn.ReadJSON(&got); err != nil {
		t.Fatalf("read update after valid subscribe: %v", err)
	}
}

func TestServeNoGoroutineLeakOnDisconnect(t *testing.T) {
	fixed := clock.FixedClock{Time: time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)}
	s := New(fixed, 5*time.Millisecond)
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	before := runtime.NumGoroutine()

	conn, _, err := websocket.DefaultDialer.Dial(wsURL(ts), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	if err := conn.WriteJSON(SubscribeRequest{Lat: 10, Lon: 10}); err != nil {
		t.Fatalf("write subscribe: %v", err)
	}
	var got UpdateResponse
	if err := conn.ReadJSON(&got); err != nil {
		t.Fatalf("read update: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= before {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Errorf("goroutine leak: before=%d after=%d", before, runtime.NumGoroutine())
}
