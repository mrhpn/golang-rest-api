package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/config"
	mw "github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/modules/health"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"github.com/mrhpn/go-rest-api/internal/types"
	"gorm.io/gorm"
)

// can't put db and cfg to context?
func Register(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	api := router.Group("/api")

	// ----------------------- Set up (Wiring) ----------------------- //
	// health
	healthH := health.NewHandler()

	// users
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)
	userH := users.NewHandler(userService)

	// ----------------------- ROUTES ----------------------- //

	// ----------------------- health ----------------------- //
	h := api.Group("/health")
	{
		h.GET("/", healthH.Check)
	}

	// ----------------------- users ----------------------- //
	u := api.Group("/users")
	u.Use(mw.RequireAuth(cfg.JWTSecret))
	{
		u.POST("", mw.AllowRoles(types.RoleSuperAdmin, types.RoleAdmin), userH.Create)
		u.DELETE("/:id", mw.AllowRoles(types.RoleSuperAdmin, types.RoleAdmin), userH.Delete)
		u.PUT("/:id/restore", mw.AllowRoles(types.RoleSuperAdmin, types.RoleAdmin), userH.Restore)
		u.GET("/:id", userH.Get)
	}
}
