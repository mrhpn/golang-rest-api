package main

import (
	"net/http"

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
		// limit total request size to 8MB
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, ctx.Cfg.MaxRequestBodySize)
		c.Next()
	})

	// global middlewares
	router.Use(middlewares.Recovery())
	router.Use(middlewares.CORS(ctx))
	router.Use(middlewares.RequestID(ctx.Cfg.AppEnv))
	router.Use(middlewares.RequestLogger())

	// register all module routes
	routes.Register(router, ctx)

	// ram management
	router.MaxMultipartMemory = 8 << 20

	return router
}
