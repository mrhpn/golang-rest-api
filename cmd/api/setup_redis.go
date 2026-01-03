package main

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/database"
)

func setupRedis(cfg *config.Config) *redis.Client {
	if !cfg.Redis.Enabled {
		log.Info().Msg("Redis is disabled, skipping connection")
		return nil
	}

	log.Info().
		Str("host", cfg.Redis.Host).
		Str("port", cfg.Redis.Port).
		Bool("enabled", cfg.Redis.Enabled).
		Msg("Connecting to Redis...")

	client, err := database.ConnectRedis(&cfg.Redis)
	if err != nil {
		log.Fatal().Msg("failed to connect Redis.")
	}

	log.Info().Msg("✅ Redis connected successfully")
	return client
}
