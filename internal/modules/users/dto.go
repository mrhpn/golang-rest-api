package users

import (
	"time"

	"github.com/mrhpn/go-rest-api/internal/security"
)

type IDParam struct {
	ID string `uri:"id" binding:"required,ulid"`
}

type CreateUserRequest struct {
	Email    string        `json:"email" binding:"required,email"`
	Password string        `json:"password" binding:"required,min=8"`
	Role     security.Role `json:"role" binding:"required,oneof=admin employee user"`
}

// UserResponse returns necessary data about a user
type UserResponse struct {
	ID        string              `json:"id"`
	Email     string              `json:"email"`
	Role      security.Role       `json:"role"`
	Status    security.UserStatus `json:"status"`
	CreatedAt string              `json:"created_at"`
}

// ToUserResponse converts a User model to UserResponse DTO
func ToUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}

// ToUserResponseList converts a slice of User models to UserResponse DTOs
func ToUserResponseList(users []*User) []UserResponse {
	if users == nil {
		return []UserResponse{}
	}
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}
