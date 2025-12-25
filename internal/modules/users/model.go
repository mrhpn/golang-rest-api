package users

import (
	"github.com/mrhpn/go-rest-api/internal/models"
	"github.com/mrhpn/go-rest-api/internal/security"
)

type User struct {
	models.BaseModel

	Email        string        `gorm:"uniqueIndex;not null"`
	PasswordHash string        `gorm:"not null"`
	Role         security.Role `gorm:"type:varchar(20);not null;default:'user'"`
}
