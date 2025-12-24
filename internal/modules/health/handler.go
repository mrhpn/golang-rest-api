package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/httpx"
)

// func Register(r *gin.RouterGroup) {
// 	r.GET("/health", func(c *gin.Context) {
// 		httpx.OK(c, http.StatusOK, gin.H{
// 			"status": "healthy",
// 		})
// 	})
// }

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Check(c *gin.Context) {
	httpx.OK(c, http.StatusOK, gin.H{"status": "healthy"})
}
