package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/httpx"
)

type Handler struct {
	service Service
}

func NewHandler(r *gin.RouterGroup, service Service) {
	h := &Handler{service: service}

	users := r.Group("/users")
	{
		users.POST("", h.create)
		users.GET("/:id", h.get)
		users.DELETE("/:id", h.delete)
		users.PUT("/:id/restore", h.restore)
	}
}

func (h *Handler) create(c *gin.Context) {
	var req CreateUserRequest

	if err := httpx.BindJSON(c, &req); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.service.Create(c, req)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusCreated, UserResponse{ID: user.ID, Email: user.Email})
}

func (h *Handler) get(c *gin.Context) {
	var params IDParam

	if err := httpx.BindURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.service.GetById(c, params.ID)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(
		c,
		http.StatusOK,
		UserResponse{ID: user.ID, Email: user.Email},
	)
}

func (h *Handler) delete(c *gin.Context) {
	var params IDParam

	if err := httpx.BindURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err := h.service.Delete(c, params.ID); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) restore(c *gin.Context) {
	var params IDParam

	if err := httpx.BindURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err := h.service.Restore(c, params.ID); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, nil)
}
