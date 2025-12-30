package database

import (
	"context"
	"fmt"
	"time"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// ConnectRedis establishes a connection to Redis
func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	if !cfg.Enabled {
		return nil, nil // Redis is optional
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		// Connection pool settings
		PoolSize:     10,
		MinIdleConns: 5,
		// Timeouts
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		// Connection lifecycle
		ConnMaxIdleTime: 5 * time.Minute,
		ConnMaxLifetime: 30 * time.Minute,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info().
		Str("host", cfg.Host).
		Str("port", cfg.Port).
		Int("db", cfg.DB).
		Msg("Redis connection established successfully")

	return rdb, nil
}
