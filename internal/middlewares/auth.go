package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/types"
	"github.com/rs/zerolog/log"
)

type UserClaims struct {
	UserID string     `json:"user_id"`
	Role   types.Role `json:"role"`
	jwt.RegisteredClaims
}

type contextKey string

const userKey contextKey = "user_identity"

// RequireAuth validates the JWT and injects claims into the context
func RequireAuth(ctx *app.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check for Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing token"})
			return
		}

		// parse and validate JWT
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &UserClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return []byte(ctx.Cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid token"})
			return
		}

		// tag the logger with UserID for better traceability
		l := log.Ctx(c.Request.Context()).
			With().
			Str("user_id", claims.UserID).
			Str("role", string(claims.Role)).
			Logger()

		ctx := context.WithValue(c.Request.Context(), userKey, claims)
		c.Request = c.Request.WithContext(l.WithContext(ctx))

		c.Next()
	}
}

// AllowRoles is a middleware factory for RBAC
func AllowRoles(allowedRoles ...types.Role) gin.HandlerFunc {
	for _, role := range allowedRoles {
		if !types.ValidRoles[role] {
			panic(fmt.Sprintf("invalid role '%s' passed to AllowRoles middleware. check your route definitions!", role))
		}
	}

	return func(c *gin.Context) {
		// pull claims from context
		val := c.Request.Context().Value(userKey)
		claims, ok := val.(*UserClaims)
		if !ok || claims == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "identity not found in context"})
			return
		}

		// check if user's role is in allowed list
		isAllowed := slices.Contains(allowedRoles, claims.Role)

		if !isAllowed {
			required := make([]string, len(allowedRoles))
			for i, r := range allowedRoles {
				required[i] = string(r)
			}
			log.Ctx(c.Request.Context()).Warn().
				Str("user_id", claims.UserID).
				Str("role", string(claims.Role)).
				Str("required_roles", strings.Join(required, ",")).
				Msg("access denied due to insufficient role")

			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
			return
		}

		c.Next()
	}
}

// GetUser is a helper for services to grab the current user
func GetUser(ctx context.Context) (*UserClaims, error) {
	claims, ok := ctx.Value(userKey).(*UserClaims)
	if !ok || claims == nil {
		return nil, fmt.Errorf("user identity not found in context")
	}
	return claims, nil
}
