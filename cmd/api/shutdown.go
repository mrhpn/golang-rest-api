package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/mrhpn/go-rest-api/internal/config"
)

const (
	timeoutSecond   = 5 * time.Second
	shutdownTimeout = 15 * time.Second
)

func gracefulShutdown(cfg *config.Config, srv *http.Server, db *gorm.DB) {
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
	log.Info().Dur("timeout", shutdownTimeout).Msg("Shutting down HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("❌ HTTP server shutdown failed")
	} else {
		log.Info().Msg("✓ HTTP server shut down gracefully")
	}

	// ensure database connection is closed on exit
	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("❌ Failed to get database connection")
	} else if sqlDB != nil {
		// Close with timeout
		closeCtx, closeCancel := context.WithTimeout(context.Background(), timeoutSecond)
		defer closeCancel()

		closed := make(chan error, 1)
		go func() {
			closed <- sqlDB.Close()
		}()

		select {
		case err = <-closed:
			if err != nil {
				log.Error().Err(err).Msg("❌ Failed to close database connection")
			} else {
				log.Info().Msg("✓ Database connection closed")
			}
		case <-closeCtx.Done():
			log.Warn().Msg("⚠ Database connection close timeout")
		}
	}

	log.Info().Msg("✓ Server exited gracefully")
}
