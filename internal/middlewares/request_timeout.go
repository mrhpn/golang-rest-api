package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/constants"
)

// RequestTimeout creates a middleware that cancels the request context after the specified duration.
// This prevents long-running requests from consuming resources indefinitely.
// Important notes:
//   - DON'T USE this middleware for websocket routes
//   - DON'T write to gin.Context from goroutines spawned within handlers
//   - The timeout applies to the entire request handler chain
//   - When timeout occurs, the context is cancelled, propagating to DB queries and external calls
func RequestTimeout(timeoutDuration time.Duration) gin.HandlerFunc {
	if timeoutDuration <= 0 {
		timeoutDuration = constants.RequestTimeoutSecond * time.Second
	}

	return timeout.New(
		timeout.WithTimeout(timeoutDuration),
		timeout.WithResponse(createTimeoutResponseHandler()),
	)
}

func createTimeoutResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Ctx(c.Request.Context()).Warn().Msg("request timeout exceeded")

		c.JSON(http.StatusGatewayTimeout, gin.H{
			"success": false,
			"error": gin.H{
				"code":    apperror.ErrRequestTimeout.Code,
				"message": apperror.ErrRequestTimeout.Message,
			},
		})
	}
}
