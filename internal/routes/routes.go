package routes

import (
	"github.com/gin-gonic/gin"

	_ "github.com/mrhpn/go-rest-api/docs" // integrates docs
	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/constants"
	"github.com/mrhpn/go-rest-api/internal/modules/auth"
	"github.com/mrhpn/go-rest-api/internal/modules/health"
	"github.com/mrhpn/go-rest-api/internal/modules/media"
	"github.com/mrhpn/go-rest-api/internal/modules/posts"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
)

// Register registers app's api endpoints
func Register(router *gin.Engine, appCtx *app.Context) {
	// API versioning: v1 is the current version
	api := router.Group(constants.APIVersionPrefix)

	// --- repositories --- //
	userR := users.NewRepository(appCtx.DB)
	postR := posts.NewRepository(appCtx.DB)

	// --- services --- //
	userS := users.NewService(userR)
	postS := posts.NewService(postR)
	authS := auth.NewService(userS, appCtx.SecurityHandler)

	// --- handlers --- //
	authH := auth.NewHandler(authS, appCtx)
	userH := users.NewHandler(userS)
	postH := posts.NewHandler(postS)
	mediaH := media.NewHandler(appCtx.MediaService)
	healthH := health.NewHandler(appCtx)

	// --- routes --- //
	registerSwagger(router)
	registerHealth(router, api, healthH)

	registerAuth(api, appCtx, authH)
	registerUsers(api, appCtx, userH)
	registerMedia(api, appCtx, mediaH)
	registerPosts(api, appCtx, postH)

	registerFallbacks(router)
}
