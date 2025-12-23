package httpx

import (
	"context"

	"github.com/gin-gonic/gin"
)

func ReqCtx(c *gin.Context) context.Context {
	return c.Request.Context()
}
