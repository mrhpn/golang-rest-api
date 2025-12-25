package users

import "github.com/mrhpn/go-rest-api/internal/security"

type IDParam struct {
	ID string `uri:"id" binding:"required,ulid"`
}

type CreateUserRequest struct {
	Email    string        `json:"email" binding:"required,email"`
	Password string        `json:"password" binding:"required,min=8"`
	Role     security.Role `json:"role" binding:"required,oneof=admin employee user"`
}

type UserResponse struct {
	ID    string        `json:"id"`
	Email string        `json:"email"`
	Role  security.Role `json:"role"`
}
