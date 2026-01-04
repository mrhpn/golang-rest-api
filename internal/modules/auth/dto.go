package auth

import (
	"github.com/mrhpn/go-rest-api/internal/modules/users"
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
	ID     string              `json:"id"`
	Email  string              `json:"email"`
	Role   security.Role       `json:"role"`
	Status security.UserStatus `json:"status"`
}

// RefreshTokenResponse constructs refresh token response structure.
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// ToLoginUserResponse converts a User model to LoginUserResponse DTO
func ToLoginUserResponse(user *users.User) LoginUserResponse {
	return LoginUserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Status: user.Status,
	}
}

// ToLoginResponse converts access token and user to LoginResponse DTO
func ToLoginResponse(accessToken string, user *users.User) LoginResponse {
	return LoginResponse{
		AccessToken: accessToken,
		User:        ToLoginUserResponse(user),
	}
}

// ToRefreshTokenResponse converts access token to RefreshTokenResponse DTO
func ToRefreshTokenResponse(newAccessToken string) RefreshTokenResponse {
	return RefreshTokenResponse{
		AccessToken: newAccessToken,
	}
}
