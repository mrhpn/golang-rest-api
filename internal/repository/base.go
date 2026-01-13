package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/mrhpn/go-rest-api/internal/database"
)

type Base struct {
	DBInstance *gorm.DB
}

// DB returns a gorm.DB instance tied to the provided context.
func (r *Base) DB(ctx context.Context) *gorm.DB {
	// 1. check if there is an active transaction in the context
	if tx := database.GetTx(ctx); tx != nil {
		return tx
	}
	// 2. fallback to standard DB with context
	return r.DBInstance.WithContext(ctx)
}
