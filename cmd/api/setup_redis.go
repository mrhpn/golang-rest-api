package main

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/database"
)

func setupRedis(cfg *config.Config) (*redis.Client, func(), error) {
	if !cfg.Redis.Enabled {
		log.Info().Msg("Redis is disabled, skipping connection")
		return nil, func() {}, nil
	}

	client, err := database.ConnectRedis(&cfg.Redis)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect Redis: %w", err)
	}

	log.Info().
		Str("host", cfg.Redis.Host).
		Str("port", cfg.Redis.Port).
		Bool("enabled", cfg.Redis.Enabled).
		Msg("✅ Redis connected successfully")

	cleanup := func() {
		if closeErr := client.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("failed to close Redis connection")
		} else {
			log.Info().Msg("✓ Redis connection closed")
		}
	}

	return client, cleanup, nil
}
