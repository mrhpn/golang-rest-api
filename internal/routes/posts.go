package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/app"
	mw "github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/modules/posts"
)

func registerPosts(api *gin.RouterGroup, appCtx *app.Context, postH *posts.Handler) {
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
