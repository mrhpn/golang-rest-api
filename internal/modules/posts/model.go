package posts

import (
	"github.com/mrhpn/go-rest-api/internal/model"
)

// PostStatus represents the status of posts
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

// Post represents the db model for post
type Post struct {
	model.Base

	UserID  string     `gorm:"column:user_id;type:char(26);not null;index" json:"user_id"`
	Title   string     `gorm:"not null" json:"title"`
	Content string     `gorm:"type:text;not null" json:"content"`
	Status  PostStatus `gorm:"type:varchar(20);not null;default:'draft';index" json:"status"`
}

// TableName specifies the table name for the Post model
func (Post) TableName() string {
	return "posts"
}

// IsValidPostStatus reports whether the given status is supported by the system.
func IsValidPostStatus(status PostStatus) bool {
	switch status {
	case PostStatusDraft, PostStatusPublished, PostStatusArchived:
		return true
	default:
		return false
	}
}
