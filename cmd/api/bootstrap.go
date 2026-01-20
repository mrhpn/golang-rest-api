package main

import (
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/app"
)

func runApplication() error {
	// Setup config
	cfg := setupConfig()

	// Setup logger
	logger := setupLogger(cfg)

	// Setup database connection
	db, dbCleanup, dbErr := setupDatabase(cfg)
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

	// Setup media storage (minio)
	media, mediaCleanup, mediaErr := setupMedia(cfg)
	if mediaErr != nil {
		log.Error().Err(mediaErr).Msg("media setup failed")
		return mediaErr
	}
	defer mediaCleanup()

	appCtx := setupAppContext(cfg, db, redis, logger, media) // app context

	// Run development-only cleanup of old rate-limit keys
	app.CleanupOldRateLimitKeysOnStartup(appCtx)

	router := setupRouter(appCtx)          // router
	server := setupHTTPServer(cfg, router) // server

	// Start server and handle shutdown
	return gracefulShutdown(cfg, server)
}
