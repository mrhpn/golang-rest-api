package routes

import (
	"github.com/gin-gonic/gin"

	_ "github.com/mrhpn/go-rest-api/docs" // integrates docs
	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/constants"
)

// Register registers app's api endpoints
func Register(router *gin.Engine, appCtx *app.Context) {
	// API versioning: v1 is the current version
	api := router.Group(constants.APIVersionPrefix)

	registerSwagger(router)
	registerHealth(router, api, appCtx)
	registerAuth(api, appCtx)
	registerUsers(api, appCtx)
	registerMedia(api, appCtx)
	registerPosts(api, appCtx)
	registerFallbacks(router)
}
