package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/constants"
)

func RequestBodyLimit(maxRequestBodySize int64) gin.HandlerFunc {
	if maxRequestBodySize <= 0 {
		maxRequestBodySize = constants.RequestMaxBodySizeMB
	}

	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestBodySize)
		c.Next()
	}
}
