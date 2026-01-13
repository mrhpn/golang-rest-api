// Package model provides basic model structure for application's business models
package model

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Base includes common fields that are required in every business models
type Base struct {
	ID        string         `gorm:"primaryKey;type:char(26)" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate generates ulid before creating a database record
func (b *Base) BeforeCreate(_ *gorm.DB) error {
	if b.ID == "" {
		b.ID = ulid.Make().String()
	}
	return nil
}
