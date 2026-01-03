package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/httpx"
)

// RateLimitResult contains the result of a rate limit check
type RateLimitResult struct {
	Allowed   bool
	Remaining int
	ResetAt   time.Time
}

// RedisRateLimiter implements rate limiting using Redis with sliding window algorithm
// This solves:
// 1. Multi-instance isolation (shared state in Redis)
// 2. Fixed window bursting (sliding window prevents edge bursts)
// 3. Proper rate limit headers (including X-RateLimit-Reset)
type RedisRateLimiter struct {
	client *redis.Client
	rate   int
	window time.Duration
}

// NewRedisRateLimiter creates a new Redis-based rate limiter
func NewRedisRateLimiter(client *redis.Client, rate int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		rate:   rate,
		window: window,
	}
}

// Allow checks if a request is allowed and returns rate limit information
// Uses sliding window algorithm to prevent burst attacks at window boundaries
func (rl *RedisRateLimiter) Allow(ctx context.Context, key string) (*RateLimitResult, error) {
	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Redis key for this identifier
	redisKey := fmt.Sprintf("ratelimit:%s", key)

	// Use Redis pipeline for atomic operations
	pipe := rl.client.Pipeline()

	// Remove old entries (outside the sliding window)
	pipe.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(windowStart.Unix(), 10))

	// Count current requests in the window
	countCmd := pipe.ZCard(ctx, redisKey)

	// Add current request with timestamp as score
	pipe.ZAdd(ctx, redisKey, redis.Z{
		Score:  float64(now.Unix()),
		Member: strconv.FormatInt(now.UnixNano(), 10), // Unique member
	})

	// Set expiration on the key (cleanup)
	pipe.Expire(ctx, redisKey, rl.window)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("redis rate limit check failed: %w", err)
	}

	// Get count after adding current request
	count := int(countCmd.Val())

	// Calculate remaining requests
	remaining := max(rl.rate-count, 0)

	// Check if allowed
	allowed := count <= rl.rate

	// Calculate reset time (oldest entry in window + window duration)
	var resetAt time.Time
	if count > 0 {
		// Get the oldest entry's timestamp
		oldestCmd := rl.client.ZRangeWithScores(ctx, redisKey, 0, 0)
		if len(oldestCmd.Val()) > 0 {
			oldestTimestamp := int64(oldestCmd.Val()[0].Score)
			resetAt = time.Unix(oldestTimestamp, 0).Add(rl.window)
		} else {
			resetAt = now.Add(rl.window)
		}
	} else {
		resetAt = now.Add(rl.window)
	}

	return &RateLimitResult{
		Allowed:   allowed,
		Remaining: remaining,
		ResetAt:   resetAt,
	}, nil
}

// createRateLimitHandler is the shared implementation for Redis-based rate limiting
func createRateLimitHandler(ctx *app.Context, rate int, window time.Duration) gin.HandlerFunc {
	if !ctx.Cfg.RateLimit.Enabled {
		// Rate limiting disabled, return no-op middleware
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Validate parameters
	if rate <= 0 {
		rate = 100
	}
	if window <= 0 {
		window = time.Minute
	}

	// Check if Redis is enabled
	if ctx.Cfg.Redis.Enabled && ctx.Redis != nil {
		limiter := NewRedisRateLimiter(ctx.Redis, rate, window)

		return func(c *gin.Context) {
			// Use IP address as the rate limit key
			key := c.ClientIP()

			// Check rate limit
			result, err := limiter.Allow(c.Request.Context(), key)
			if err != nil {
				// If Redis fails, log error but allow request (fail open)
				// In production, you might want to fail closed instead
				log.Ctx(c.Request.Context()).Error().
					Err(err).
					Str("ip", key).
					Msg("rate limit check failed, allowing request")
				c.Next()
				return
			}

			// Set rate limit headers (RFC 7231 compliant)
			c.Header("X-RateLimit-Limit", strconv.Itoa(limiter.rate))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))

			if !result.Allowed {
				log.Ctx(c.Request.Context()).Warn().
					Str("ip", key).
					Str("path", c.Request.URL.Path).
					Int("rate", rate).
					Int("reset_at", int(result.ResetAt.Unix())).
					Msg("rate limit exceeded")

				httpx.Fail(
					c,
					http.StatusTooManyRequests,
					"RATE_LIMIT_EXCEEDED",
					fmt.Sprintf("rate limit exceeded. retry after %d", result.ResetAt.Unix()),
					nil,
				)
				c.Abort()
				return
			}

			c.Next()
		}
	}

	// Redis not enabled, fall back to in-memory rate limiter
	log.Warn().Msg("Redis not enabled, falling back to in-memory rate limiter (not suitable for multi-instance deployments)")
	return RateLimit(rate, window)
}

// RateLimitRedis middleware using Redis for distributed rate limiting
// Solves all three issues:
// 1. Multi-instance: Shared Redis state across all replicas
// 2. Sliding window: Prevents burst attacks at window boundaries
// 3. Proper headers: Includes X-RateLimit-Reset with Unix timestamp
func RateLimitRedis(ctx *app.Context) gin.HandlerFunc {
	rate := ctx.Cfg.RateLimit.Rate
	if rate <= 0 {
		rate = 100
	}

	window := time.Duration(ctx.Cfg.RateLimit.Window) * time.Second
	if window <= 0 {
		window = time.Minute
	}

	return createRateLimitHandler(ctx, rate, window)
}

// RateLimitRedisWithConfig creates a Redis-based rate limiter with custom rate and window
// Use this for route-specific rate limiting (e.g., stricter limits for auth endpoints)
func RateLimitRedisWithConfig(ctx *app.Context, rate int, window time.Duration) gin.HandlerFunc {
	return createRateLimitHandler(ctx, rate, window)
}
