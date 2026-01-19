package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/app"
	mw "github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"github.com/mrhpn/go-rest-api/internal/security"
)

func registerUsers(api *gin.RouterGroup, appCtx *app.Context) {
	// ----------------------- 1. Set up (Wiring) ----------------------- //
	userR := users.NewRepository(appCtx.DB)
	userS := users.NewService(userR)
	userH := users.NewHandler(userS)

	// ----------------------- 2. ROUTES ----------------------- //
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
}
