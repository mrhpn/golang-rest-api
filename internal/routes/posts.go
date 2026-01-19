package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/app"
	mw "github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/modules/posts"
)

func registerPosts(api *gin.RouterGroup, appCtx *app.Context) {
	// ----------------------- 1. Set up (Wiring) ----------------------- //
	postR := posts.NewRepository(appCtx.DB)
	postS := posts.NewService(postR)
	postH := posts.NewHandler(postS)

	// ----------------------- 2. ROUTES ----------------------- //
	postsGroup := api.Group("/posts")
	postsGroup.Use(mw.RequireAuth(appCtx))
	{
		postsGroup.POST("", postH.Create)
		postsGroup.GET("", postH.List)
		postsGroup.GET("/my", postH.ListMyPosts)
		postsGroup.GET("/:id", postH.Get)
		postsGroup.PUT("/:id", postH.Update)
		postsGroup.DELETE("/:id", postH.Delete)
	}
}
