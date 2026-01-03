// Package routes define entire application's api endpoints
package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/mrhpn/go-rest-api/docs" // integrates docs
	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/constants"
	mw "github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/modules/auth"
	"github.com/mrhpn/go-rest-api/internal/modules/health"
	"github.com/mrhpn/go-rest-api/internal/modules/media"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"github.com/mrhpn/go-rest-api/internal/security"
)

// Register registers app's api endpoints
func Register(router *gin.Engine, ctx *app.Context) {
	authRateLimitCount := ctx.Cfg.RateLimit.AuthRate
	if authRateLimitCount <= 0 {
		authRateLimitCount = constants.RateLimitAuth // Default: 5 requests per minute for auth endpoints
	}

	// API versioning: v1 is the current version
	api := router.Group(constants.APIVersionPrefix)

	// Swagger Route
	// Access at: http://localhost:8080/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.PersistAuthorization(true)))

	// ----------------------- Set up (Wiring) ----------------------- //
	// health
	healthH := health.NewHandler(ctx)

	// users
	userR := users.NewRepository(ctx.DB)
	userS := users.NewService(userR)
	userH := users.NewHandler(userS)

	// auth
	securityH := security.NewJWTHandler(
		ctx.Cfg.JWT.Secret,
		ctx.Cfg.JWT.AccessTokenExpirationSecond,
		ctx.Cfg.JWT.RefreshTokenExpirationSecond,
	)
	authS := auth.NewService(userS, securityH)
	authH := auth.NewHandler(authS, ctx)

	// media
	mediaH := media.NewHandler(ctx.MediaService)

	// ----------------------- ROUTES ----------------------- //

	// ----------------------- health ----------------------- //
	// Health endpoints (outside /api for Kubernetes/Docker health checks)
	router.GET("/health", healthH.Check)
	router.GET("/health/live", healthH.Liveness)
	router.GET("/health/ready", healthH.Readiness)

	// Also expose under /api/v1 for consistency
	api.GET("/health", healthH.Check)
	api.GET("/health/live", healthH.Liveness)
	api.GET("/health/ready", healthH.Readiness)

	// ----------------------- auth ----------------------- //
	// Apply stricter rate limiting for auth endpoints using Redis
	authGroup := api.Group("/" + constants.APIAuthPrefix)
	authGroup.Use(mw.RateLimitRedisWithConfig(ctx, authRateLimitCount, time.Minute))
	{
		authGroup.POST("/login", authH.Login)
		authGroup.POST("/refresh", authH.Refresh)
	}

	// ----------------------- users ----------------------- //
	usersGroup := api.Group("/users")
	usersGroup.Use(mw.RequireAuth(ctx))
	{
		usersGroup.GET("/:id", userH.Get)
		usersGroup.POST("", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Create)
		usersGroup.DELETE("/:id", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Delete)
		usersGroup.PUT("/:id/restore", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Restore)
	}

	// ----------------------- media ----------------------- //
	mediaGroup := api.Group("/media")
	mediaGroup.Use(mw.RequireAuth(ctx))
	{
		mediaGroup.POST("/upload/profile", mediaH.UploadProfilePicture)
	}
}
