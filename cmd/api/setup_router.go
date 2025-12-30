package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/routes"
)

func setupRouter(ctx *app.AppContext) *gin.Engine {
	httpx.RegisterValidators()
	if ctx.Cfg.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// global request body limit middleware
	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, ctx.Cfg.Http.MaxRequestBodySize)
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
	requestTimeout := time.Duration(ctx.Cfg.Http.RequestTimeoutSecond) * time.Second
	if requestTimeout <= 0 {
		requestTimeout = 30 * time.Second
	}
	router.Use(middlewares.Timeout(requestTimeout))

	// Rate limiting middleware (if enabled)
	if ctx.Cfg.RateLimit.Enabled {
		rateLimitWindow := time.Duration(ctx.Cfg.RateLimit.Window) * time.Second
		if rateLimitWindow <= 0 {
			rateLimitWindow = time.Minute
		}
		router.Use(middlewares.RateLimit(ctx.Cfg.RateLimit.Rate, rateLimitWindow))
	}

	// register all module routes
	routes.Register(router, ctx)

	// ram management
	router.MaxMultipartMemory = 8 << 20

	return router
}
