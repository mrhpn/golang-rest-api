package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/config"
)

const (
	shutdownTimeout = 15 * time.Second
)

func gracefulShutdown(cfg *config.Config, srv *http.Server) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info().
			Str("port", cfg.Port).
			Str("env", cfg.AppEnv).
			Msg("🚀 HTTP server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("❌ HTTP server failed")
		}
	}()

	// wait for shutdown signal
	<-ctx.Done()
	log.Info().Msg("Shutdown signal received, starting graceful shutdown...")

	// create context with timeout for graceful shutdown
	// stops accepting new requests, waits for existing requests to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// attempt graceful shutdown
	log.Info().Dur("timeout", shutdownTimeout).Msg("Shutting down HTTP server & other services...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("❌ HTTP server shutdown failed")
		return err
	}

	log.Info().Msg("✓ HTTP server shut down gracefully")
	log.Info().Msg("✓ Server exited gracefully")
	return nil
}
