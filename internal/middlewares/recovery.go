package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/httpx"
)

// Recovery gives ability to recover from internal server errors
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
					apperror.ErrInternal.Code,
					apperror.ErrInternal.Message,
					nil,
				)
			}
		}()
		c.Next()
	}
}
