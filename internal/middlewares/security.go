package middlewares

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy (formerly Feature-Policy)
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Strict Transport Security (HSTS) - only set in production with HTTPS
		// Uncomment and configure if using HTTPS:
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		c.Next()
	}
}
