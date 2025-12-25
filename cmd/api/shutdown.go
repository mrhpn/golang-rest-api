package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func gracefulShutdown(cfg *config.Config, srv *http.Server, db *gorm.DB) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info().Msgf("🚀 HTTP server started on port %s (env=%s) ", cfg.Port, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("❌ HTTP server failed")
		}
	}()

	// wait for shutdown signal
	<-ctx.Done()
	log.Info().Msg("1. Shutdown singnal received, starting graceful shutdown...")

	// create context with timeout for graceful shutdown
	// stops accepting new requests, waits for existing requests to finish for 10s
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// attempt graceful shutdown
	log.Info().Msg("2. Shutting down HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("❌ HTTP server shutdown failed")
	} else {
		log.Info().Msg("3. HTTP server shut down gracefully")
	}

	// ensure database connection is closed on exit
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			log.Error().Err(err).Msg("❌ Failed to close database connection")
		} else {
			log.Info().Msg("4. Database connection closed")
		}
	}
	log.Info().Msg("5. Server exited gracefully. Bye!")
}
