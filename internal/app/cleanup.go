package app

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/constants"
)

const (
	healthCheckTimeout       = 5 * time.Second
	countPerScan       int64 = 100
)

// CleanupOldRateLimitKeysOnStartup cleans up old rate limit keys stored in redis on app startup (ON DEV Env)
func CleanupOldRateLimitKeysOnStartup(ctx *Context) {
	if ctx.Redis == nil {
		return
	}

	cleanupCtx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()

	// Find all rate limit keys
	var cursor uint64
	var total int
	pattern := constants.RateLimitKeyPrefix + "*"
	for {
		keys, cur, err := ctx.Redis.Scan(cleanupCtx, cursor, pattern, countPerScan).Result()
		if err != nil {
			log.Warn().Err(err).Msg("failed to scan rate limit keys for cleanup")
			return
		}

		if len(keys) > 0 {
			// Delete all old rate limit keys
			if err = ctx.Redis.Del(cleanupCtx, keys...).Err(); err != nil {
				log.Warn().Err(err).Msg("failed to delete old rate limit keys during cleanup")
			} else {
				total += len(keys)
			}
		}
		cursor = cur
		if cursor == 0 {
			break
		}
	}
	if total > 0 {
		log.Info().Int("count", total).Msg("Cleaned up old rate limit keys on startup")
	} else {
		log.Info().Msg("🙌 Redis — No old rate limit keys deleted (Env - Dev)")
	}
}
