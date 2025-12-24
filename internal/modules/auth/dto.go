package auth

import "github.com/mrhpn/go-rest-api/internal/modules/users"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string             `json:"token"`
	User  users.UserResponse `json:"user"`
}
