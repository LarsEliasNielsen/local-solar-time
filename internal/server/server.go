package server

import (
	"net/http"

	"local-solar-time/internal/clock"
)

type SubscribeRequest struct {
	Lat float64
	Lon float64
}

type Server struct {
	Clock clock.Clock
}

func New(c clock.Clock) *Server {
	panic("not implemented")
}

func (s *Server) Handler() http.Handler {
	panic("not implemented")
}
