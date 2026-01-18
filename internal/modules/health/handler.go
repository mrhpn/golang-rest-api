// Package health exposes health and readiness checks for the application.
package health

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/security"
)

const healthCheckTimeout = 5 * time.Second

// Handler handles application health and readiness check HTTP endpoints.
type Handler struct {
	appCtx *app.Context
}

// NewHandler constructs a health Handler with access to the application context.
func NewHandler(appCtx *app.Context) *Handler {
	return &Handler{appCtx: appCtx}
}

// Check health godoc
//
//	@Summary		Check health
//	@Description	Check health status of server (liveness probe)
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	health.Response
//	@Router			/health [get]
func (h *Handler) Check(c *gin.Context) {
	httpx.OK(c, http.StatusOK, ToResponse("healthy"))
}

// Readiness checks if the service is ready to accept traffic
//
//	@Summary		Check readiness
//	@Description	Check if service is ready to accept traffic (readiness probe)
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	health.Response
//	@Failure		503	{object}	health.Response
//	@Router			/health/ready [get]
func (h *Handler) Readiness(c *gin.Context) {
	checks := make(map[string]string)
	allHealthy := true

	ctx, cancel := context.WithTimeout(httpx.ReqCtx(c), healthCheckTimeout)
	defer cancel()

	// Check database connectivity
	if !checkDB(ctx, h.appCtx, checks) {
		allHealthy = false
	}

	// Check MinIO storage connectivity
	if !checkStorage(ctx, h.appCtx, checks) {
		allHealthy = false
	}

	// Check Redis connectivity (if enabled)
	if !checkRedis(ctx, h.appCtx, checks) {
		allHealthy = false
	}

	status := "ready"
	httpStatus := http.StatusOK
	if !allHealthy {
		status = "not_ready"
		httpStatus = http.StatusServiceUnavailable
	}

	httpx.OK(c, httpStatus, ToResponse(status, checks))
}

// Liveness checks if the service is alive
//
//	@Summary		Check liveness
//	@Description	Check if service is alive (liveness probe)
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	health.Response
//	@Router			/health/live [get]
func (h *Handler) Liveness(c *gin.Context) {
	httpx.OK(c, http.StatusOK, ToResponse("alive"))
}

// RateLimitStatus checks the status of the rate limit
//
//	@Summary		Check rate limit status
//	@Description	Check the status of the rate limit
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	health.RateLimitResponse
//	@Router			/health/rate-limit/status [get]
func (h *Handler) RateLimitStatus(c *gin.Context) {
	if !h.appCtx.Cfg.RateLimit.Enabled {
		httpx.OK(c, http.StatusOK, ToResponse("disabled", map[string]string{"rate_limit": "disabled"}))
		return
	}

	if h.appCtx.Redis == nil {
		httpx.FailWithError(c, apperror.New(
			apperror.Internal,
			apperror.ErrInternal.Code,
			"redis enabled for ratelimit, but redis client not initialized",
		))
		return
	}

	ctx, cancel := context.WithTimeout(httpx.ReqCtx(c), healthCheckTimeout)
	defer cancel()

	keys, err := h.appCtx.Redis.Keys(ctx, "ratelimit:*").Result()
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	// Get details for each key
	details := make(map[string]RateLimitKeysDetails)
	for _, key := range keys {
		valStr, _ := h.appCtx.Redis.Get(ctx, key).Result()
		ttl, _ := h.appCtx.Redis.TTL(ctx, key).Result()

		count, _ := strconv.Atoi(valStr)

		details[key] = RateLimitKeysDetails{
			Count: count,
			TTL:   ttl.String(),
		}
	}

	httpx.OK(c, http.StatusOK, ToRedisRateLimitResponse(len(keys), details))
}

// ResetRateLimit resets the rate limit
//
//	@Summary		Reset rate limit
//	@Description	Reset the rate limit
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	health.RedisRateLimitResetResponse
//	@Router			/health/rate-limit/reset [post]
func (h *Handler) ResetRateLimit(c *gin.Context) {
	if h.appCtx.Cfg.AppEnv != "development" {
		httpx.Fail(
			c,
			http.StatusForbidden,
			security.ErrForbiddenInProd.Code,
			security.ErrForbiddenInProd.Message,
			nil,
		)
		return
	}

	if h.appCtx.Redis == nil {
		httpx.FailWithError(c, apperror.New(
			apperror.Internal,
			apperror.ErrInternal.Code,
			"redis enabled for ratelimit, but redis client not initialized",
		))
		return
	}

	ctx, cancel := context.WithTimeout(httpx.ReqCtx(c), healthCheckTimeout)
	defer cancel()

	keys, err := h.appCtx.Redis.Keys(ctx, "ratelimit:*").Result()
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if len(keys) > 0 {
		if err := h.appCtx.Redis.Del(ctx, keys...).Err(); err != nil {
			httpx.FailWithError(c, err)
			return
		}
	}

	httpx.OK(c, http.StatusOK, ToRedisRateLimitResetResponse("rate limits reset", len(keys)))
}
