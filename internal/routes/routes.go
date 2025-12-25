package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/app"
	mw "github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/modules/auth"
	"github.com/mrhpn/go-rest-api/internal/modules/health"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"github.com/mrhpn/go-rest-api/internal/security"
)

func Register(router *gin.Engine, ctx *app.AppContext) {
	api := router.Group("/api")

	// ----------------------- Set up (Wiring) ----------------------- //
	// health
	healthH := health.NewHandler()

	// users
	userRepo := users.NewRepository(ctx.DB)
	userService := users.NewService(userRepo)
	userH := users.NewHandler(userService)

	// auth
	securityHandler := security.NewJWTHandler(
		ctx.Cfg.JWT.Secret,
		ctx.Cfg.JWT.AccessTokenExpirationSecond,
		ctx.Cfg.JWT.RefreshTokenExpirationSecond,
	)
	authService := auth.NewService(userService, securityHandler)
	authH := auth.NewHandler(authService)

	// ----------------------- ROUTES ----------------------- //

	// ----------------------- health ----------------------- //
	h := api.Group("/health")
	{
		h.GET("/", healthH.Check)
	}

	// ----------------------- auth ----------------------- //
	api.POST("/login", authH.Login)
	api.POST("/auth/refresh", authH.Refresh)

	// ----------------------- users ----------------------- //
	u := api.Group("/users")
	u.Use(mw.RequireAuth(ctx))
	{
		u.GET("/:id", userH.Get)
		u.POST("", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Create)
		u.DELETE("/:id", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Delete)
		u.PUT("/:id/restore", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Restore)
	}
}
