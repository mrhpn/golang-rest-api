package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/modules/auth"
	"github.com/mrhpn/go-rest-api/internal/security"
)

type contextKey string

const userKey contextKey = "user_identity"

// RequireAuth validates the JWT and injects claims into the context
func RequireAuth(ctx *app.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. check for Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			httpx.Fail(
				c,
				http.StatusUnauthorized,
				auth.ErrUnauthorized.Code,
				auth.ErrUnauthorized.Message,
				nil,
			)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. parse and validate JWT
		claims, err := ctx.SecurityHandler.ValidateToken(tokenString)
		if err != nil {
			// Map security errors to proper HTTP responses
			httpx.FailWithError(c, err)
			return
		}

		// 3. tag the logger with UserID for better traceability
		l := log.Ctx(c.Request.Context()).
			With().
			Str("user_id", claims.UserID).
			Str("role", string(claims.Role)).
			Logger()

		// 4. inject claims into req context
		reqCtx := context.WithValue(c.Request.Context(), userKey, claims)
		c.Request = c.Request.WithContext(l.WithContext(reqCtx))

		c.Next()
	}
}

// AllowRoles is a middleware factory for RBAC
func AllowRoles(allowedRoles ...security.Role) gin.HandlerFunc {
	// validation check on startup
	for _, role := range allowedRoles {
		if !security.IsValidRole(role) {
			panic(fmt.Sprintf("invalid role '%s' passed to AllowRoles middleware. check your route definitions!", role))
		}
	}

	return func(c *gin.Context) {
		// 1. pull claims from context
		val := c.Request.Context().Value(userKey)
		claims, ok := val.(*security.UserClaims)
		if !ok || claims == nil {
			httpx.Fail(
				c,
				http.StatusInternalServerError,
				auth.ErrIdentityNotFoundInContext.Code,
				auth.ErrIdentityNotFoundInContext.Message,
				nil,
			)
			return
		}

		// 2. check if user's role is in allowed list
		if !slices.Contains(allowedRoles, claims.Role) {
			log.Ctx(c.Request.Context()).Warn().
				Str("user_id", claims.UserID).
				Str("role", string(claims.Role)).
				Interface("required_roles", allowedRoles).
				Msg("access denied due to insufficient role")

			httpx.Fail(
				c,
				http.StatusForbidden,
				auth.ErrForbidden.Code,
				auth.ErrForbidden.Message,
				nil,
			)
			return
		}

		c.Next()
	}
}

// GetUser is a helper for services to grab the current user
func GetUser(ctx context.Context) (*security.UserClaims, error) {
	claims, ok := ctx.Value(userKey).(*security.UserClaims)
	if !ok || claims == nil {
		return nil, errors.New("user identity not found in context")
	}
	return claims, nil
}
