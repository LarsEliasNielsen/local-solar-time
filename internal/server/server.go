package server

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"

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
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func New(c clock.Clock, cadence time.Duration) *Server {
	return &Server{Clock: c, Cadence: cadence}
}

// Handler upgrades each request to a WebSocket connection and serves it on its own goroutine.
func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		go s.serve(conn)
	})
}

// serve reads the client's one-time subscribe message, then pushes a freshly
// computed solar payload on every tick until the connection closes.
func (s *Server) serve(conn *websocket.Conn) {
	defer func() { _ = conn.Close() }()

	var req SubscribeRequest
	if err := conn.ReadJSON(&req); err != nil {
		return
	}

	ticker := time.NewTicker(s.Cadence)
	defer ticker.Stop()

	for range ticker.C {
		result := solar.Compute(s.Clock.Now(), req.Lat, req.Lon)
		if err := conn.WriteJSON(result); err != nil {
			return
		}
	}
}
