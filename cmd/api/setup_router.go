package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/constants"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/routes"
)

func setupRouter(ctx *app.Context) *gin.Engine {
	httpx.RegisterValidators()
	if ctx.Cfg.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// global request body limit middleware
	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, ctx.Cfg.HTTP.MaxRequestBodySize)
		c.Next()
	})

	// Security headers (should be early in the chain)
	router.Use(middlewares.SecurityHeaders())

	// Recovery middleware (should be early to catch panics)
	router.Use(middlewares.Recovery())

	// CORS middleware
	router.Use(middlewares.CORS(ctx))

	// Request ID middleware (for tracing)
	router.Use(middlewares.RequestID(ctx.Cfg.AppEnv))

	// Request logger middleware
	router.Use(middlewares.RequestLogger())

	// Request timeout middleware
	requestTimeout := time.Duration(ctx.Cfg.HTTP.RequestTimeoutSecond) * time.Second
	if requestTimeout <= 0 {
		requestTimeout = constants.RequestTimeoutSecond * time.Second
	}
	router.Use(middlewares.RequestTimeout(requestTimeout))

	// Rate limiting middleware (if enabled)
	// Uses Redis if available, falls back to in-memory if not
	if ctx.Cfg.RateLimit.Enabled {
		router.Use(middlewares.RateLimitRedis(ctx))
	}

	// register all module routes
	routes.Register(router, ctx)

	// ram management
	router.MaxMultipartMemory = constants.MaxMultipartMemory

	return router
}
