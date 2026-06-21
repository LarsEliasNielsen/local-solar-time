package config

import (
	"io"

	"github.com/rs/zerolog"
)

// NewLogger returns a structured JSON logger writing to w, defaulting to
// Info level so Debug-only fields (e.g. coordinates) are suppressed unless
// explicitly raised.
func NewLogger(w io.Writer) zerolog.Logger {
	return zerolog.New(w).Level(zerolog.InfoLevel).With().Timestamp().Logger()
}
