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

	tokenPair, user, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User: users.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		}})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshTokenRequest
	if err := httpx.BindJSON(c, &req); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	ctx := httpx.ReqCtx(c)

	newAccessToken, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, RefreshTokenResponse{
		AccessToken: newAccessToken,
	})
}
