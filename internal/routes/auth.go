package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/constants"
	mw "github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/modules/auth"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
)

func registerAuth(api *gin.RouterGroup, appCtx *app.Context) {
	// Apply stricter rate limiting for auth endpoints using Redis
	authRateLimit := appCtx.Cfg.RateLimit.AuthRate
	if authRateLimit == "" {
		authRateLimit = constants.RateLimitAuth // Default: 7 requests per minute for auth endpoints
	}

	// ----------------------- 1. Set up (Wiring) ----------------------- //
	// user
	userR := users.NewRepository(appCtx.DB)
	userS := users.NewService(userR)

	// auth
	authS := auth.NewService(userS, appCtx.SecurityHandler)
	authH := auth.NewHandler(authS, appCtx)

	// ----------------------- 2. ROUTES ----------------------- //
	authGroup := api.Group("/" + constants.APIAuthPrefix)
	authGroup.Use(mw.RateLimitRedisWithConfig(appCtx, authRateLimit))
	{
		authGroup.POST("/login", authH.Login)
		authGroup.POST("/refresh", authH.Refresh)
	}
}
