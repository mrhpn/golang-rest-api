package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mrhpn/go-rest-api/internal/config"
)

const (
	poolSize           = 10
	minIdleConns       = 5
	dialTimeout        = 5 * time.Second
	readTimeout        = 3 * time.Second
	writeTimeout       = 3 * time.Second
	connMaxIdleTime    = 5 * time.Minute
	connMaxLifetime    = 30 * time.Minute
	healthCheckTimeout = 5 * time.Second
)

// ConnectRedis establishes a connection to Redis
func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		// Connection pool settings
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
		// Timeouts
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		// Connection lifecycle
		ConnMaxIdleTime: connMaxIdleTime,
		ConnMaxLifetime: connMaxLifetime,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}
