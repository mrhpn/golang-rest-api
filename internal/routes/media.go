package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/app"
	mw "github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/modules/media"
)

func registerMedia(api *gin.RouterGroup, appCtx *app.Context) {
	// ----------------------- 1. Set up (Wiring) ----------------------- //
	mediaH := media.NewHandler(appCtx.MediaService)

	// ----------------------- 2. ROUTES ----------------------- //
	mediaGroup := api.Group("/media")
	mediaGroup.Use(mw.RequireAuth(appCtx))
	{
		mediaGroup.POST("/upload/profile", mediaH.UploadProfilePicture)
	}
}
