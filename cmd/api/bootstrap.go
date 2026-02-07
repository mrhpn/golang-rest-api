package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/constants"
)

func runApplication() error {
	// Setup config
	cfg, cfgErr := setupConfig()
	if cfgErr != nil {
		log.Error().Err(cfgErr).Msg("config setup failed")
		return cfgErr
	}

	// Setup logger
	logger := setupLogger(cfg)

	// Setup database connection & materics logging
	dbMetricsCtx, dbMetricsCancel := context.WithCancel(context.Background())
	defer dbMetricsCancel()

	db, dbCleanup, dbErr := setupDatabase(dbMetricsCtx, cfg)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("database setup failed")
		return dbErr
	}
	defer dbCleanup()

	// Setup redis (optional)
	redis, redisCleanup, redisErr := setupRedis(cfg)
	if redisErr != nil {
		log.Error().Err(redisErr).Msg("redis setup failed")
		return redisErr
	}
	defer redisCleanup()

	// Setup mediaSvc storage (minio)
	mediaSvc, mediaCleanup, mediaErr := setupMedia(cfg)
	if mediaErr != nil {
		log.Error().Err(mediaErr).Msg("media setup failed")
		return mediaErr
	}
	defer mediaCleanup()

	appCtx := setupAppContext(cfg, db, redis, logger, mediaSvc) // app context

	// Run development-only cleanup of old rate-limit keys
	if cfg.AppEnv == constants.EnvDev {
		app.CleanupOldRateLimitKeysOnStartup(appCtx)
	}

	router := setupRouter(appCtx)          // router
	server := setupHTTPServer(cfg, router) // server

	// Start server and handle shutdown
	return gracefulShutdown(cfg, server)
}
