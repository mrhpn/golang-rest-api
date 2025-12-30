package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/modules/auth"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Log panic with full context
				log.Ctx(c.Request.Context()).
					Error().
					Interface("panic", r).
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("panic recovered")

				httpx.Fail(
					c,
					http.StatusInternalServerError,
					auth.ErrInternal.Code,
					auth.ErrInternal.Message,
					nil,
				)
			}
		}()
		c.Next()
	}
}
