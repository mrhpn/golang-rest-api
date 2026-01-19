package posts

import (
	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"github.com/mrhpn/go-rest-api/internal/security"
	"github.com/mrhpn/go-rest-api/internal/timex"
)

type IDParam struct {
	ID string `uri:"id" binding:"required,ulid"`
}

type CreatePostRequest struct {
	Title   string     `json:"title" binding:"required,min=1,max=200"`
	Content string     `json:"content" binding:"required,min=1"`
	Status  PostStatus `json:"status" binding:"omitempty,oneof=draft published archived"`
}

type UpdatePostRequest struct {
	Title   string     `json:"title" binding:"omitempty,min=1,max=200"`
	Content string     `json:"content" binding:"omitempty,min=1"`
	Status  PostStatus `json:"status" binding:"omitempty,oneof=draft published archived"`
}

// PostResponse returns necessary data about a post
type PostResponse struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Status    PostStatus     `json:"status"`
	Author    AuthorResponse `json:"author"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

// AuthorResponse returns minimum necessary data about author of a post
type AuthorResponse struct {
	ID    string        `json:"id"`
	Email string        `json:"email"`
	Role  security.Role `json:"role"`
}

// ToPostResponse converts a Post model to PostResponse DTO
func ToPostResponse(post *Post) PostResponse {
	return PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Status:    post.Status,
		Author:    ToAuthorResponse(&post.User),
		CreatedAt: timex.ToAPIDateTimeFormat(post.CreatedAt),
		UpdatedAt: timex.ToAPIDateTimeFormat(post.UpdatedAt),
	}
}

// ToAuthorResponse converts a User model to mini AuthorResponse DTO
func ToAuthorResponse(user *users.User) AuthorResponse {
	return AuthorResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
}

// ToPostResponseList converts a slice of Post models to PostResponse DTOs
func ToPostResponseList(posts []*Post) []PostResponse {
	if posts == nil {
		return []PostResponse{}
	}
	responses := make([]PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = ToPostResponse(post)
	}
	return responses
}
