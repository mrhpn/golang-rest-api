// Package health exposes health and readiness checks for the application.
package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/httpx"
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
//	@Success		200	{object}	health.HealthResponse
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
//	@Success		200	{object}	health.HealthResponse
//	@Failure		503	{object}	health.HealthResponse
//	@Router			/health/ready [get]
func (h *Handler) Readiness(c *gin.Context) {
	checks := make(map[string]string)
	allHealthy := true

	ctx, cancel := context.WithTimeout(c.Request.Context(), healthCheckTimeout)
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
//	@Success		200	{object}	health.HealthResponse
//	@Router			/health/live [get]
func (h *Handler) Liveness(c *gin.Context) {
	httpx.OK(c, http.StatusOK, ToResponse("alive"))
}
