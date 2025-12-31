package auth

import "github.com/mrhpn/go-rest-api/internal/modules/users"

// LoginRequest constructs login request structure.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse constructs login response structure. This will not include RefreshToken since it is sent via http-only cookie for security.
type LoginResponse struct {
	AccessToken string             `json:"access_token"`
	User        users.UserResponse `json:"user"`
}

// RefreshTokenResponse constructs refresh token response structure.
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}
