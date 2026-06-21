package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"

	"local-solar-time/internal/clock"
	"local-solar-time/internal/solar"
)

type SubscribeRequest struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Server struct {
	Clock   clock.Clock
	Cadence time.Duration
	Logger  zerolog.Logger

	done      chan struct{}
	closeOnce sync.Once
	wg        sync.WaitGroup
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func New(c clock.Clock, cadence time.Duration) *Server {
	return &Server{Clock: c, Cadence: cadence, Logger: zerolog.Nop(), done: make(chan struct{})}
}

// Handler upgrades each request to a WebSocket connection and serves it on its own goroutine.
func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.serve(conn)
		}()
	})
}

// Shutdown signals every in-flight connection to close - each gets a
// server-initiated close frame, mirroring the courtesy already extended to
// client-initiated closes - then blocks until all of them have exited.
func (s *Server) Shutdown() {
	s.closeOnce.Do(func() { close(s.done) })
	s.wg.Wait()
}

// serve reads subscribe messages until one passes validation - replying with
// an ErrorResponse and retrying on failure - then pushes a freshly computed
// solar payload on every tick until the connection closes.
func (s *Server) serve(conn *websocket.Conn) {
	defer func() { _ = conn.Close() }()

	var req SubscribeRequest
	for {
		if err := conn.ReadJSON(&req); err != nil {
			return
		}
		err := validate(req)
		if err == nil {
			break
		}
		if writeErr := conn.WriteJSON(ErrorResponse{Error: err.Error()}); writeErr != nil {
			return
		}
	}

	s.Logger.Debug().Float64("lat", req.Lat).Float64("lon", req.Lon).Msg("client subscribed")

	// Read continuously in the background so gorilla/websocket can process
	// control frames - in particular, so a client-initiated close frame is
	// seen and answered, completing the closing handshake cleanly instead
	// of leaving the client to time out with an abnormal closure.
	closed := make(chan struct{})
	go func() {
		defer close(closed)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	ticker := time.NewTicker(s.Cadence)
	defer ticker.Stop()

	for {
		select {
		case <-closed:
			return
		case <-s.done:
			closeMsg := websocket.FormatCloseMessage(websocket.CloseGoingAway, "server shutting down")
			_ = conn.WriteMessage(websocket.CloseMessage, closeMsg)
			return
		case <-ticker.C:
			result := solar.Compute(s.Clock.Now(), req.Lat, req.Lon)
			if err := conn.WriteJSON(toUpdateResponse(result)); err != nil {
				return
			}
		}
	}
}
