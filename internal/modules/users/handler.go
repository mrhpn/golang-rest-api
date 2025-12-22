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
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, httpx.ErrBadRequest, err.Error())
		return
	}

	user, err := h.service.Create(req)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.ErrInternal, err.Error())
		return
	}

	httpx.OK(c, http.StatusCreated, UserResponse{ID: user.ID, Email: user.Email})
}

func (h *Handler) get(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetById(id)
	if err != nil {
		httpx.Fail(c, http.StatusNotFound, httpx.ErrNotFound, "user not found")
		return
	}

	httpx.OK(c, http.StatusOK, UserResponse{ID: user.ID, Email: user.Email})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.ErrInternal, err.Error())
		return
	}

	httpx.OK(c, http.StatusNoContent, nil)
}

func (h *Handler) restore(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Restore(id); err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.ErrInternal, err.Error())
		return
	}

	httpx.OK(c, http.StatusOK, nil)
}
