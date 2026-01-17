package httpx

import (
	"context"

	"github.com/gin-gonic/gin"
)

// ReqCtx returns the request context from the Gin context.
func ReqCtx(c *gin.Context) context.Context {
	return c.Request.Context()
}
