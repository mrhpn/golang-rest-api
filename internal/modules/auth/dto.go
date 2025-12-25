package auth

import "github.com/mrhpn/go-rest-api/internal/modules/users"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Note: LoginResponse will not include RefreshToken since it is sent via http-only cookie for security
type LoginResponse struct {
	AccessToken string             `json:"access_token"`
	User        users.UserResponse `json:"user"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}
