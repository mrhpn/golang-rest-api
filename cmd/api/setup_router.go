package main

import (
	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/constants"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/routes"
)

func setupRouter(ctx *app.Context) *gin.Engine {
	httpx.RegisterValidators()

	if ctx.Cfg.AppEnv != constants.EnvDev {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 🛡️ Recovery middleware
	router.Use(middlewares.Recovery())

	// 📏 Request body limit middleware
	router.Use(middlewares.RequestBodyLimit(ctx.Cfg.HTTP.MaxRequestBodySize))

	// 🔒 Security headers
	router.Use(middlewares.SecurityHeaders())

	// 🌐 CORS middleware
	router.Use(middlewares.CORS(ctx.Cfg.HTTP.AllowedOrigins))

	// 🆔 Request ID middleware (for tracing)
	router.Use(middlewares.RequestID(ctx.Cfg.AppEnv))

	// ⏱️ Request timeout middleware
	router.Use(middlewares.RequestTimeout(ctx.Cfg.HTTP.RequestTimeoutSecond))

	// 📝 Request logger middleware
	router.Use(middlewares.RequestLogger())

	// 🚦 Rate limiting middleware (if enabled)
	// Uses Redis if available, falls back to in-memory if not
	if ctx.Cfg.RateLimit.Enabled {
		router.Use(middlewares.RateLimitRedis(ctx))
	}

	// Register all module routes
	routes.Register(router, ctx)

	// ram management
	router.MaxMultipartMemory = constants.MaxMultipartMemory

	return router
}
