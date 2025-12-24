package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/app"
)

func CORS(ctx *app.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowOrigin := ""

		// determine if the origin is allowed
		if len(ctx.Cfg.AllowedOrigins) > 0 && ctx.Cfg.AllowedOrigins[0] == "*" {
			allowOrigin = origin // !use origin, instead of "*"
		} else if slices.Contains(ctx.Cfg.AllowedOrigins, origin) {
			allowOrigin = origin
		}

		// set headers if allowed
		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Crendentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		// handle preflight
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
