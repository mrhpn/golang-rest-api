package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/httpx"
)

const refreshTokenCookieName = "refresh_token"

// Handler handles authentication-related HTTP endpoints such as login, token refresh, and access control–protected actions.
type Handler struct {
	authService Service
	ctx         *app.Context
}

// NewHandler constructs an authentication Handler with its required dependencies
func NewHandler(authService Service, ctx *app.Context) *Handler {
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
//	@Param			request	body		auth.LoginRequest	true	"Login Credentials"
//	@Success		200		{object}	auth.LoginResponse
//	@Failure		400		{object}	httpx.ErrorResponse
//	@Failure		401		{object}	httpx.ErrorResponse
//	@Failure		500		{object}	httpx.ErrorResponse
//	@Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := httpx.BindAndValidateJSON(c, &req); err != nil {
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
		refreshTokenCookieName,            // name
		tokenPair.RefreshToken,            // value
		cookieMaxAge,                      // max age
		"/",                               // path
		"",                                // domain (empty for localhost)
		h.ctx.Cfg.AppEnv != "development", // secure (set to true in production with HTTPS)
		true,                              // httpOnly
	)

	httpx.OK(c, http.StatusOK, ToLoginResponse(tokenPair.AccessToken, user))
}

// Refresh token godoc
//
//	@Summary		Refresh token
//	@Description	Generate a new access token, request must have valid refresh token in cookie
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	auth.RefreshTokenResponse
//	@Failure		400	{object}	httpx.ErrorResponse
//	@Failure		401	{object}	httpx.ErrorResponse
//	@Failure		500	{object}	httpx.ErrorResponse
//	@Router			/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	// read form cookie instead of json body
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		httpx.Fail(
			c,
			http.StatusUnauthorized,
			errRefreshTokenMissing.Code,
			errRefreshTokenMissing.Message,
			nil,
		)
		return
	}

	ctx := httpx.ReqCtx(c)

	newAccessToken, err := h.authService.RefreshToken(ctx, refreshToken)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, ToRefreshTokenResponse(newAccessToken))
}
