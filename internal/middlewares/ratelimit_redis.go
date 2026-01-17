package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/ulule/limiter/v3"
	ginlimit "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/constants"
	"github.com/mrhpn/go-rest-api/internal/httpx"
)

// createRateLimitHandler creates a rate limit handler using ulule/limiter
// rateStr should be in ulule/limiter format: "100-M" (100 per minute), "50-H" (50 per hour), "10-S" (10 per second)
func createRateLimitHandler(ctx *app.Context, rateStr string) gin.HandlerFunc {
	if !ctx.Cfg.RateLimit.Enabled {
		// Rate limiting disabled, return no-op middleware
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Validate and set default if empty
	if rateStr == "" {
		rateStr = constants.RateLimit
	}

	// Parse the rate limit
	rateLimit, err := limiter.NewRateFromFormatted(rateStr)
	if err != nil {
		log.Error().
			Err(err).
			Str("rate", rateStr).
			Msg("failed to parse rate limit, using default")
		rateLimit, _ = limiter.NewRateFromFormatted(constants.RateLimit)
	}

	var store limiter.Store

	// Check if Redis is enabled
	if ctx.Cfg.Redis.Enabled && ctx.Redis != nil {
		// Use Redis store for distributed rate limiting
		var redisStore limiter.Store
		redisStore, err = redisstore.NewStoreWithOptions(ctx.Redis, limiter.StoreOptions{
			Prefix: "ratelimit",
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to create Redis rate limit store, falling back to memory")
			// Fall back to memory store
			store = memory.NewStore()
		} else {
			store = redisStore
		}
	} else {
		// Use in-memory store (not suitable for multi-instance deployments)
		log.Warn().Msg("Redis not enabled, using in-memory rate limiter (not suitable for multi-instance deployments)")
		store = memory.NewStore()
	}

	// Create limiter instance
	instance := limiter.New(store, rateLimit)

	// Create Gin middleware
	ginMiddleware := ginlimit.NewMiddleware(instance)

	// Wrap with custom error handling to match error format
	return func(c *gin.Context) {
		// Call the ulule/limiter middleware
		ginMiddleware(c)

		// Check if the request was rate limited (ulule/limiter sets status 429)
		if c.Writer.Status() == http.StatusTooManyRequests {
			// Get rate limit info from headers (ulule/limiter sets these)
			key := c.ClientIP()
			limit := c.GetHeader("X-RateLimit-Limit")
			reset := c.GetHeader("X-RateLimit-Reset")

			log.Ctx(httpx.ReqCtx(c)).Warn().
				Str("ip", key).
				Str("path", c.Request.URL.Path).
				Str("rate", rateStr).
				Str("limit", limit).
				Str("reset", reset).
				Msg("rate limit exceeded")

			// Override the response to match your error format
			httpx.Fail(
				c,
				http.StatusTooManyRequests,
				apperror.ErrTooManyRequests.Code,
				fmt.Sprintf("rate limit exceeded. retry after %s", reset),
				nil,
			)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitRedis middleware using ulule/limiter with Redis backend
// This is the global rate limiter that automatically skips OPTIONS requests and auth routes
// (auth routes have their own stricter rate limiter)
func RateLimitRedis(ctx *app.Context) gin.HandlerFunc {
	rateStr := ctx.Cfg.RateLimit.Rate
	if rateStr == "" {
		rateStr = constants.RateLimit
	}

	// Get the base rate limit handler
	baseHandler := createRateLimitHandler(ctx, rateStr)

	// Wrap it with skip logic for OPTIONS requests and auth routes
	return func(c *gin.Context) {
		// Skip rate limiting for OPTIONS (CORS preflight) requests
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// Skip global rate limiting for auth routes (they have their own stricter limits)
		if strings.HasPrefix(c.Request.URL.Path, constants.APIAuthPath) {
			c.Next()
			return
		}

		// Apply rate limiting for all other routes
		baseHandler(c)
	}
}

// RateLimitRedisWithConfig creates a rate limiter with custom rate string
// rateStr should be in ulule/limiter format: "100-M" (100 per minute), "50-H" (50 per hour), "10-S" (10 per second)
// Use this for route-specific rate limiting (e.g., stricter limits for auth endpoints)
func RateLimitRedisWithConfig(ctx *app.Context, rateStr string) gin.HandlerFunc {
	return createRateLimitHandler(ctx, rateStr)
}
