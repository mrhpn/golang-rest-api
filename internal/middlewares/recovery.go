package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/httpx"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Interface("panic", r).Msg("panic recovered")
				httpx.Fail(
					c,
					http.StatusInternalServerError,
					"INTERNAL",
					"internal server error",
					nil,
				)
				c.Abort()
			}
		}()
		c.Next()
	}
}
