package users

import "github.com/mrhpn/go-rest-api/internal/security"

type iDParam struct {
	ID string `uri:"id" binding:"required,ulid"`
}

type createUserRequest struct {
	Email    string        `json:"email" binding:"required,email"`
	Password string        `json:"password" binding:"required,min=8"`
	Role     security.Role `json:"role" binding:"required,oneof=admin employee user"`
}

// UserResponse returns necessary data about a user
type UserResponse struct {
	ID    string        `json:"id"`
	Email string        `json:"email"`
	Role  security.Role `json:"role"`
}
