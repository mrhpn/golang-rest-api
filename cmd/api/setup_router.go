package main

import (
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

	// global middlewares
	router.Use(middlewares.Recovery())
	router.Use(middlewares.CORS(ctx))
	router.Use(middlewares.RequestID(ctx.Cfg.AppEnv))
	router.Use(middlewares.RequestLogger())

	// register all module routes
	routes.Register(router, ctx)

	return router
}
