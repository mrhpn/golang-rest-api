package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/httpx"
)

// RequestID middleware sets request id to the context, into the response header, and logs for better traceability.
func RequestID(env string) gin.HandlerFunc {
	switch env {
	case "development", "production", "testing":
	default:
		env = "development"
	}

	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 1. store in gin context for manual retrieval (c.GetString("RequestID"))
		c.Set("RequestID", requestID)
		// 2. set in response header for the frontend/client to consume
		c.Header("X-Request-ID", requestID)

		// 3. attach request_id to logger
		l := log.With().
			Str("env", env).
			Str("request_id", requestID).
			Logger()
		// 4. inject this logger into the Stanard Library Context
		c.Request = c.Request.WithContext(l.WithContext(httpx.ReqCtx(c)))

		c.Next()
	}
}
