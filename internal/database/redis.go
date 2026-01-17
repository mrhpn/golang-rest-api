package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/constants"
)

// ConnectRedis establishes a connection to Redis
func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		// Connection pool settings
		PoolSize:     constants.RedisPoolSize,
		MinIdleConns: constants.RedisMinIdleConns,
		// Timeouts
		DialTimeout:  constants.RedisDialTimeoutSecond,
		ReadTimeout:  constants.RedisReadTimeoutSecond,
		WriteTimeout: constants.RedisWriteTimeoutSecond,
		// Connection lifecycle
		ConnMaxIdleTime: constants.RedisConnMaxIdleTimeMinute,
		ConnMaxLifetime: constants.RedisConnMaxLifetimeMinute,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), constants.RedisHealthCheckTimeoutSecond)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return rdb, nil
}
