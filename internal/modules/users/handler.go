package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/httpx"
)

type Handler struct {
	userService Service
}

func NewHandler(userService Service) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := httpx.BindJSON(c, &req); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.userService.Create(httpx.ReqCtx(c), req)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(
		c,
		http.StatusCreated,
		UserResponse{ID: user.ID, Email: user.Email, Role: user.Role},
	)
}

func (h *Handler) Get(c *gin.Context) {
	var params IDParam

	if err := httpx.BindURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.userService.GetById(httpx.ReqCtx(c), params.ID)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(
		c,
		http.StatusOK,
		UserResponse{ID: user.ID, Email: user.Email, Role: user.Role},
	)
}

func (h *Handler) Delete(c *gin.Context) {
	var params IDParam

	if err := httpx.BindURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err := h.userService.Delete(httpx.ReqCtx(c), params.ID); err != nil {
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

	if err := h.userService.Restore(httpx.ReqCtx(c), params.ID); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, nil)
}
