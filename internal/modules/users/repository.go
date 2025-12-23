package users

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindById(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	SoftDelete(ctx context.Context, id string) (int64, error)
	Restore(ctx context.Context, id string) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) FindById(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err // db issue
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
		return nil, err // db issue
	}

	return &user, nil
}

func (r *repository) SoftDelete(ctx context.Context, id string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&User{}, "id = ?", id)
	return result.RowsAffected, result.Error
}

func (r *repository) Restore(ctx context.Context, id string) (int64, error) {
	result := r.db.Unscoped().
		WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Update("deleted_at", nil)

	return result.RowsAffected, result.Error
}
