package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
)

type Handler struct {
	authService Service
}

func NewHandler(authService Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := httpx.BindJSON(c, &req); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	ctx := httpx.ReqCtx(c)

	token, user, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, LoginResponse{Token: token, User: users.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}})
}
