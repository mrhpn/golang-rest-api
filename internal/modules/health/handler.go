// Package health exposes health and readiness checks for the application.
package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/rs/zerolog/log"
)

// Handler handles application health and readiness check HTTP endpoints.
type Handler struct {
	appCtx *app.Context
}

// NewHandler constructs a health Handler with access to the application context.
func NewHandler(appCtx *app.Context) *Handler {
	return &Handler{appCtx: appCtx}
}

// healthResponse represents the health check response
type healthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// Check health godoc
//
//	@Summary		Check health
//	@Description	Check health status of server (liveness probe)
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	healthResponse
//	@Router			/health [get]
func (h *Handler) Check(c *gin.Context) {
	httpx.OK(c, http.StatusOK, healthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Readiness checks if the service is ready to accept traffic
//
//	@Summary		Check readiness
//	@Description	Check if service is ready to accept traffic (readiness probe)
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	healthResponse
//	@Failure		503	{object}	healthResponse
//	@Router			/health/ready [get]
func (h *Handler) Readiness(c *gin.Context) {
	checks := make(map[string]string)
	allHealthy := true

	// Check database connectivity
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	sqlDB, err := h.appCtx.DB.DB()
	if err != nil {
		checks["db"] = "unhealthy: failed to get database connection"
		allHealthy = false
	} else {
		if err := sqlDB.PingContext(ctx); err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("database health check failed")
			checks["db"] = "unhealthy: " + err.Error()
			allHealthy = false
		} else {
			stats := sqlDB.Stats()
			checks["db"] = "healthy"
			checks["db_open_conns"] = fmt.Sprintf("%d", stats.OpenConnections)
			checks["db_idle_conns"] = fmt.Sprintf("%d", stats.Idle)
		}
	}

	// Check MinIO storage connectivity
	if h.appCtx.MediaService != nil {
		if err := h.appCtx.MediaService.HealthCheck(ctx); err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("minio health check failed")
			checks["storage"] = "unhealthy: " + err.Error()
			allHealthy = false
		} else {
			checks["storage"] = "healthy"
		}
	} else {
		checks["storage"] = "unhealthy: media service not initialized"
		allHealthy = false
	}

	// Check Redis connectivity (if enabled)
	if h.appCtx.Cfg.Redis.Enabled {
		if h.appCtx.Redis != nil {
			if err := h.appCtx.Redis.Ping(ctx).Err(); err != nil {
				log.Ctx(ctx).Warn().Err(err).Msg("redis health check failed")
				checks["redis"] = "unhealthy: " + err.Error()
				allHealthy = false
			} else {
				checks["redis"] = "healthy"
			}
		} else {
			checks["redis"] = "unhealthy: redis client not initialized"
			allHealthy = false
		}
	} else {
		checks["redis"] = "disabled"
		// Redis is optional, so we don't mark as unhealthy if disabled
	}

	status := "ready"
	httpStatus := http.StatusOK
	if !allHealthy {
		status = "not_ready"
		httpStatus = http.StatusServiceUnavailable
	}

	httpx.OK(c, httpStatus, healthResponse{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	})
}

// Liveness checks if the service is alive
//
//	@Summary		Check liveness
//	@Description	Check if service is alive (liveness probe)
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	healthResponse
//	@Router			/health/live [get]
func (h *Handler) Liveness(c *gin.Context) {
	httpx.OK(c, http.StatusOK, healthResponse{
		Status:    "alive",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
