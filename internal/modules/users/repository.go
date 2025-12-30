package users

import (
	"context"
	"errors"

	appErr "github.com/mrhpn/go-rest-api/internal/errors"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindById(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Delete(ctx context.Context, id string) (int64, error)
	Restore(ctx context.Context, id string) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		// Wrap database errors to preserve context while maintaining client-safe messages
		return appErr.Wrap(
			appErr.Internal,
			ErrDatabaseError.Code,
			"failed to create user",
			err,
		)
	}
	return nil
}

func (r *repository) FindById(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		// Wrap database errors to preserve context
		return nil, appErr.Wrap(
			appErr.Internal,
			ErrDatabaseError.Code,
			"failed to find user",
			err,
		)
	}

	return &user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		// Wrap database errors to preserve context
		return nil, appErr.Wrap(
			appErr.Internal,
			ErrDatabaseError.Code,
			"failed to find user",
			err,
		)
	}

	return &user, nil
}

func (r *repository) Delete(ctx context.Context, id string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&User{}, "id = ?", id)
	if result.Error != nil {
		return 0, appErr.Wrap(
			appErr.Internal,
			ErrDatabaseError.Code,
			"failed to delete user",
			result.Error,
		)
	}
	return result.RowsAffected, nil
}

func (r *repository) Restore(ctx context.Context, id string) (int64, error) {
	result := r.db.Unscoped().
		WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Update("deleted_at", nil)

	if result.Error != nil {
		return 0, appErr.Wrap(
			appErr.Internal,
			ErrDatabaseError.Code,
			"failed to restore user",
			result.Error,
		)
	}

	return result.RowsAffected, nil
}
