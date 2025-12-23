package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/database"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/routes"
	"github.com/rs/zerolog/log"
)

func main() {
	// ----- ✅ 1. load env & configs ----- //
	_ = godotenv.Load()
	cfg := config.Load()

	// ----- ✅ 2. setup logger ----- //
	logger := app.SetupLogger(cfg.AppEnv)
	log.Logger = logger

	// ----- ✅ 3. connect to database ----- //
	db, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Database connection failed")
	}

	// ----- ✅ 4. register validators ----- //
	httpx.RegisterValidators()
	if cfg.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ----- ✅ 5. setup router and register routes ----- //
	router := gin.New()
	router.Use(middlewares.Recovery())
	router.Use(middlewares.RequestID(cfg.AppEnv))
	router.Use(middlewares.RequestLogger())
	routes.Register(router, db)

	// ----- ✅ 6. setup (start/stop) HTTP server ----- //
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// start server in a goroutine
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

	// 7 ensure database connection is closed on exit
	// ----- ✅ 7. close database last ----- //
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
