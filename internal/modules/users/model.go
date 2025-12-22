package users

import "github.com/mrhpn/go-rest-api/internal/models"

type User struct {
	models.BaseModel

	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
}
