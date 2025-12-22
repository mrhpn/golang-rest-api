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
	"github.com/mrhpn/go-rest-api/internal/modules/health"
	"github.com/rs/zerolog/log"
)

func main() {
	// ----- ✅ load env ----- //
	_ = godotenv.Load()

	// ----- ✅ load configs ----- //
	cfg := config.Load()

	// ----- ✅ connect to database ----- //
	db, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Database connection failed")
	}
	// ensure database connection is closed on exit
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		log.Info().Msg("✅ Database connection closed")
	}()

	// ----- ✅ setup logger ----- //
	logger := app.SetupLogger(cfg.AppEnv)
	log.Logger = logger

	// ----- ✅ setup router and routes ----- //
	router := gin.New()
	router.Use(gin.Recovery())

	// register health check routes
	health.Register(router)

	// ----- ✅ setup (start/stop) HTTP server ----- //
	if cfg.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// start server in a goroutine
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info().Msg("🚀 HTTP server started on port " + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("❌ HTTP server failed")
		}
	}()

	// wait for shutdown signal
	<-ctx.Done()
	log.Info().Msg("⚠️ Shutdown singnal received")

	// create context with timeout for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// attempt graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("❌ Graceful shutdown failed")
	}
	log.Info().Msg("✅ Server exited gracefully")
}
