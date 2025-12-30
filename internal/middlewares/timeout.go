package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/rs/zerolog/log"
)

// Timeout creates a middleware that cancels the request context after the specified duration
// This prevents long-running requests from consuming resources indefinitely
func Timeout(timeout time.Duration) gin.HandlerFunc {
	if timeout <= 0 {
		timeout = 30 * time.Second // default: 30 seconds
	}

	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace the request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to signal completion
		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		// Wait for either completion or timeout
		select {
		case <-done:
			// Request completed successfully
		case <-ctx.Done():
			// Timeout occurred
			if ctx.Err() == context.DeadlineExceeded {
				log.Ctx(c.Request.Context()).Warn().
					Dur("timeout", timeout).
					Str("path", c.Request.URL.Path).
					Msg("request timeout exceeded")

				// Only send response if not already sent
				if !c.Writer.Written() {
					httpx.Fail(
						c,
						http.StatusRequestTimeout,
						"REQUEST_TIMEOUT",
						"request timeout exceeded",
						nil,
					)
					c.Abort()
				}
			}
		}
	}
}
