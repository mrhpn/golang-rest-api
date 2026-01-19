package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/security"
)

func registerFallbacks(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		httpx.Fail(
			c,
			http.StatusNotFound,
			security.ErrRouteNotFound.Code,
			security.ErrRouteNotFound.Message,
			nil,
		)
	})

	router.NoMethod(func(c *gin.Context) {
		httpx.Fail(
			c,
			http.StatusMethodNotAllowed,
			security.ErrMethodNotAllowed.Code,
			security.ErrMethodNotAllowed.Message,
			nil,
		)
	})
}
