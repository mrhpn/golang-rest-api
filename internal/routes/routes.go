package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/app"
	mw "github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/modules/auth"
	"github.com/mrhpn/go-rest-api/internal/modules/health"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"github.com/mrhpn/go-rest-api/internal/types"
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
	authService := auth.NewService(userService, ctx.Cfg.JWT.Secret)
	authH := auth.NewHandler(authService)

	// ----------------------- ROUTES ----------------------- //

	// ----------------------- health ----------------------- //
	h := api.Group("/health")
	{
		h.GET("/", healthH.Check)
	}

	// ----------------------- auth ----------------------- //
	api.POST("/login", authH.Login)

	// ----------------------- users ----------------------- //
	u := api.Group("/users")
	u.Use(mw.RequireAuth(ctx))
	{
		u.POST("", mw.AllowRoles(types.RoleSuperAdmin, types.RoleAdmin), userH.Create)
		u.DELETE("/:id", mw.AllowRoles(types.RoleSuperAdmin, types.RoleAdmin), userH.Delete)
		u.PUT("/:id/restore", mw.AllowRoles(types.RoleSuperAdmin, types.RoleAdmin), userH.Restore)
		u.GET("/:id", userH.Get)
	}
}
