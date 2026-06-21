// Command local-solar-time streams apparent solar time and Sun position
// for a client-supplied location over WebSocket.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"local-solar-time/internal/clock"
	"local-solar-time/internal/config"
	"local-solar-time/internal/server"
)

const version = "0.1.0"

const (
	shutdownTimeout   = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
)

func main() {
	logger := config.NewLogger(os.Stdout)

	cfg, err := config.Load()
	if err != nil {
		logger.Error().Err(err).Msg("load config")
		os.Exit(1)
	}

	srv := server.New(clock.WallClock{}, cfg.Cadence)
	srv.Logger = logger

	addr := fmt.Sprintf(":%d", cfg.Port)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		logger.Info().Str("version", version).Int("port", cfg.Port).Msg("starting server")
		serveErr <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msg("server failed")
			os.Exit(1)
		}
	case <-ctx.Done():
		logger.Info().Msg("shutting down")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("http server shutdown")
		}
		srv.Shutdown()
		logger.Info().Msg("shutdown complete")
	}
}
