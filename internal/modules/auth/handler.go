package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
)

const RefreshTokenCookieName = "refresh_token"

type Handler struct {
	authService Service
	ctx         *app.AppContext
}

func NewHandler(authService Service, ctx *app.AppContext) *Handler {
	return &Handler{
		authService: authService,
		ctx:         ctx,
	}
}

// Login godoc
//
//	@Summary		User Login
//	@Description	Authenticates a user and returns access token in response and set refresh token in httpOnly cookie
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"Login Credentials"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	httpx.ErrorResponse
//	@Failure		401		{object}	httpx.ErrorResponse
//	@Failure		500		{object}	httpx.ErrorResponse
//	@Router			/login [post]
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

	// set refresh token in cookie
	cookieMaxAge := h.ctx.Cfg.JWT.RefreshTokenExpirationSecond
	c.SetCookie(
		RefreshTokenCookieName,            // name
		tokenPair.RefreshToken,            // value
		cookieMaxAge,                      // max age
		"/",                               // path
		"",                                // domain (empty for localhost)
		h.ctx.Cfg.AppEnv != "development", // secure (set to true in production with HTTPS)
		true,                              // httpOnly
	)

	httpx.OK(c, http.StatusOK, LoginResponse{
		AccessToken: tokenPair.AccessToken,
		User: users.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		}})
}

// Refresh token godoc
//
//	@Summary		Refresh token
//	@Description	Generate a new access token, request must have valid refresh token in cookie
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	RefreshTokenResponse
//	@Failure		400	{object}	httpx.ErrorResponse
//	@Failure		401	{object}	httpx.ErrorResponse
//	@Failure		500	{object}	httpx.ErrorResponse
//	@Router			/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	// read form cookie instead of json body
	refreshToken, err := c.Cookie(RefreshTokenCookieName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	ctx := httpx.ReqCtx(c)

	newAccessToken, err := h.authService.RefreshToken(ctx, refreshToken)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, RefreshTokenResponse{
		AccessToken: newAccessToken,
	})
}
