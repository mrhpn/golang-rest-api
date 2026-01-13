package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/constants"
)

// Timeout creates a middleware that cancels the request context after the specified duration
// This prevents long-running requests from consuming resources indefinitely
func Timeout(timeout time.Duration) gin.HandlerFunc {
	if timeout <= 0 {
		timeout = constants.RequestTimeoutSecond * time.Second
	}

	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace the request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to signal completion
		c.Next()

		// Wait for either completion or timeout
		if ctx.Err() == context.DeadlineExceeded {
			log.Ctx(c.Request.Context()).Warn().Msg("request timeout exceeded")

			if !c.Writer.Written() {
				c.AbortWithStatus(http.StatusRequestTimeout)
			}
		}
	}
}
