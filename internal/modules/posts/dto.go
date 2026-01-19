package posts

import (
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
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Status    PostStatus `json:"status"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

// ToPostResponse converts a Post model to PostResponse DTO
func ToPostResponse(post *Post) PostResponse {
	return PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		Status:    post.Status,
		CreatedAt: timex.ToAPIDateTimeFormat(post.CreatedAt),
		UpdatedAt: timex.ToAPIDateTimeFormat(post.UpdatedAt),
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
