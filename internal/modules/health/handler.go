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

// Check health godoc
//
//	@Summary		Check health
//	@Description	Check health status of server
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Returns {"status": "healthy"}"
//	@Router			/health [get]
func (h *Handler) Check(c *gin.Context) {
	httpx.OK(c, http.StatusOK, gin.H{"status": "healthy"})
}
