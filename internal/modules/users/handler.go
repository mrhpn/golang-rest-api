package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/httpx"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateUserRequest

	if err := httpx.BindJSON(c, &req); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.service.Create(httpx.ReqCtx(c), req)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusCreated, UserResponse{ID: user.ID, Email: user.Email})
}

func (h *Handler) Get(c *gin.Context) {
	var params IDParam

	if err := httpx.BindURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.service.GetById(httpx.ReqCtx(c), params.ID)
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

func (h *Handler) Delete(c *gin.Context) {
	var params IDParam

	if err := httpx.BindURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err := h.service.Delete(httpx.ReqCtx(c), params.ID); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Restore(c *gin.Context) {
	var params IDParam

	if err := httpx.BindURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err := h.service.Restore(httpx.ReqCtx(c), params.ID); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, nil)
}
