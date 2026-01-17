package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/httpx"
)

// RequestLogger logs request-related logs such as unique request_id, method, path, status etc.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Ctx(httpx.ReqCtx(c)).
			Info().
			Str("request_id", c.GetString("RequestID")).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("latency", time.Since(start)).
			Msg("request completed")
	}
}
