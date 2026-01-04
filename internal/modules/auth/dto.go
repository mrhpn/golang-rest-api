package auth

import (
	"github.com/mrhpn/go-rest-api/internal/security"
)

// LoginRequest constructs login request structure.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse constructs login response structure. This will not include RefreshToken since it is sent via http-only cookie for security.
type LoginResponse struct {
	AccessToken string            `json:"access_token"`
	User        LoginUserResponse `json:"user"`
}

type LoginUserResponse struct {
	ID    string        `json:"id"`
	Email string        `json:"email"`
	Role  security.Role `json:"role"`
}

// RefreshTokenResponse constructs refresh token response structure.
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}
