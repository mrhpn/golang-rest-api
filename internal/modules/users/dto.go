package users

type IDParam struct {
	ID string `uri:"id" binding:"required,ulid"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
