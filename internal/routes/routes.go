package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mrhpn/go-rest-api/docs"
	"github.com/mrhpn/go-rest-api/internal/app"
	mw "github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/modules/auth"
	"github.com/mrhpn/go-rest-api/internal/modules/health"
	"github.com/mrhpn/go-rest-api/internal/modules/media"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"github.com/mrhpn/go-rest-api/internal/security"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Register(router *gin.Engine, ctx *app.AppContext) {
	api := router.Group("/api")

	// Swagger Route
	// Access at: http://localhost:8080/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.PersistAuthorization(true)))

	// ----------------------- Set up (Wiring) ----------------------- //
	// health
	healthH := health.NewHandler()

	// users
	userR := users.NewRepository(ctx.DB)
	userS := users.NewService(userR)
	userH := users.NewHandler(userS)

	// auth
	securityH := security.NewJWTHandler(
		ctx.Cfg.JWT.Secret,
		ctx.Cfg.JWT.AccessTokenExpirationSecond,
		ctx.Cfg.JWT.RefreshTokenExpirationSecond,
	)
	authS := auth.NewService(userS, securityH)
	authH := auth.NewHandler(authS, ctx)

	// media
	mediaH := media.NewHandler(ctx.MediaService)

	// ----------------------- ROUTES ----------------------- //

	// ----------------------- health ----------------------- //
	api.GET("/health", healthH.Check)

	// ----------------------- auth ----------------------- //
	api.POST("/login", authH.Login)
	api.POST("/auth/refresh", authH.Refresh)

	// ----------------------- users ----------------------- //
	usersGroup := api.Group("/users")
	usersGroup.Use(mw.RequireAuth(ctx))
	{
		usersGroup.GET("/:id", userH.Get)
		usersGroup.POST("", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Create)
		usersGroup.DELETE("/:id", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Delete)
		usersGroup.PUT("/:id/restore", mw.AllowRoles(security.RoleSuperAdmin, security.RoleAdmin), userH.Restore)
	}

	// ----------------------- media ----------------------- //
	mediaGroup := api.Group("/media")
	mediaGroup.Use(mw.RequireAuth(ctx))
	{
		mediaGroup.POST("/upload/profile", mediaH.UploadProfilePicture)
		mediaGroup.POST("/upload/thumbnail", mediaH.UploadThumbnail)
	}
}
