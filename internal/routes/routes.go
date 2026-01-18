// Package routes define entire application's api endpoints
package routes

import (
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
	"github.com/mrhpn/go-rest-api/internal/modules/posts"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"github.com/mrhpn/go-rest-api/internal/security"
)

// Register registers app's api endpoints
func Register(router *gin.Engine, appCtx *app.Context) {
	authRateLimit := appCtx.Cfg.RateLimit.AuthRate
	if authRateLimit == "" {
		authRateLimit = constants.RateLimitAuth // Default: 7 requests per minute for auth endpoints
	}

	// API versioning: v1 is the current version
	api := router.Group(constants.APIVersionPrefix)

	// Swagger Route
	// Access at: http://localhost:8080/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.PersistAuthorization(true)))

	// ----------------------- Set up (Wiring) ----------------------- //
	// health
	healthH := health.NewHandler(appCtx)

	// users
	userR := users.NewRepository(appCtx.DB)
	userS := users.NewService(userR)
	userH := users.NewHandler(userS)

	// auth
	// Use SecurityHandler from context instead of creating a new one
	authS := auth.NewService(userS, appCtx.SecurityHandler)
	authH := auth.NewHandler(authS, appCtx)

	// media
	mediaH := media.NewHandler(appCtx.MediaService)

	// posts
	postR := posts.NewRepository(appCtx.DB)
	postS := posts.NewService(postR)
	postH := posts.NewHandler(postS)

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
	authGroup.Use(mw.RateLimitRedisWithConfig(appCtx, authRateLimit))
	{
		authGroup.POST("/login", authH.Login)
		authGroup.POST("/refresh", authH.Refresh)
	}

	// ----------------------- users ----------------------- //
	usersGroup := api.Group("/users")
	usersGroup.Use(mw.RequireAuth(appCtx))
	{
		usersGroup.GET("", userH.List)
		usersGroup.GET("/:id", userH.Get)
		usersGroup.POST("", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Create)
		usersGroup.DELETE("/:id", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Delete)
		usersGroup.PUT("/:id/restore", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Restore)
		usersGroup.PUT("/:id/block", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Block)
		usersGroup.PUT("/:id/reactivate", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Reactivate)
	}

	// ----------------------- media ----------------------- //
	mediaGroup := api.Group("/media")
	mediaGroup.Use(mw.RequireAuth(appCtx))
	{
		mediaGroup.POST("/upload/profile", mediaH.UploadProfilePicture)
	}

	// ----------------------- posts ----------------------- //
	postsGroup := api.Group("/posts")
	postsGroup.Use(mw.RequireAuth(appCtx))
	{
		postsGroup.POST("", postH.Create)
		postsGroup.GET("", postH.List)
		postsGroup.GET("/my", postH.ListMyPosts)
		postsGroup.GET("/:id", postH.Get)
		postsGroup.PUT("/:id", postH.Update)
		postsGroup.DELETE("/:id", postH.Delete)
	}
}
