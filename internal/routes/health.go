package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/modules/health"
)

func registerHealth(router *gin.Engine, api *gin.RouterGroup, appCtx *app.Context) {
	// ----------------------- 1. Set up (Wiring) ----------------------- //
	healthH := health.NewHandler(appCtx)

	// ----------------------- 2. ROUTES ----------------------- //
	// Health endpoints (outside /api for Kubernetes/Docker health checks)
	router.GET("/health", healthH.Check)
	router.GET("/health/live", healthH.Liveness)
	router.GET("/health/ready", healthH.Readiness)

	healthGroup := api.Group("/health")
	{
		healthGroup.GET("/", healthH.Check)
		healthGroup.GET("/live", healthH.Liveness)
		healthGroup.GET("/ready", healthH.Readiness)
		rateLimitGroup := healthGroup.Group("/rate-limit")
		{
			rateLimitGroup.GET("/status", healthH.RateLimitStatus)
			rateLimitGroup.POST("/reset", healthH.ResetRateLimit)
		}
	}
}
