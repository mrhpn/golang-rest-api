package users

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/pagination"
	repo "github.com/mrhpn/go-rest-api/internal/repository"
	"github.com/mrhpn/go-rest-api/internal/security"
)

// Repository defines the persistence operations for user entities.
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, opts *pagination.QueryOptions) ([]*User, int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Restore(ctx context.Context, id string) (int64, error)
	Block(ctx context.Context, id string) (int64, error)
	Reactivate(ctx context.Context, id string) (int64, error)
	Activate(ctx context.Context, id string) (int64, error)
}

type repository struct {
	repo.Base
}

// NewRepository constructs a users Repository backed by a GORM database.
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		Base: repo.Base{
			DBInstance: db,
		},
	}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	err := r.DB(ctx).Create(user).Error
	if err != nil {
		// Wrap database errors to preserve context while maintaining client-safe messages
		return apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to create user",
			err,
		)
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.DB(ctx).First(&user, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errUserNotFound
		}
		// Wrap database errors to preserve context
		return nil, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to find user",
			err,
		)
	}

	return &user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.DB(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errUserNotFound
		}
		// Wrap database errors to preserve context
		return nil, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to find user",
			err,
		)
	}

	return &user, nil
}

func (r *repository) List(ctx context.Context, opts *pagination.QueryOptions) ([]*User, int64, error) {
	var users []*User
	var total int64

	// 1. Get total count
	err := r.DB(ctx).Model(&User{}).
		Scopes(pagination.SearchScope(opts)).
		Count(&total).Error
	if err != nil {
		return nil, 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to count users",
			err,
		)
	}

	// 2. Fetch data
	// If we add relationships later, don't forget to use Preload here
	err = r.DB(ctx).
		Scopes(pagination.Paginate(opts)).
		Find(&users).Error
	if err != nil {
		return nil, 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to find users",
			err,
		)
	}

	return users, total, nil
}

func (r *repository) Delete(ctx context.Context, id string) (int64, error) {
	result := r.DB(ctx).Delete(&User{}, "id = ?", id)
	if result.Error != nil {
		return 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to delete user",
			result.Error,
		)
	}
	return result.RowsAffected, nil
}

func (r *repository) Restore(ctx context.Context, id string) (int64, error) {
	result := r.DB(ctx).
		Unscoped().
		Model(&User{}).
		Where("id = ?", id).
		Update("deleted_at", nil)

	if result.Error != nil {
		return 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to restore user",
			result.Error,
		)
	}

	return result.RowsAffected, nil
}

func (r *repository) Block(ctx context.Context, id string) (int64, error) {
	result := r.DB(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Update("status", security.UserStatusBlocked)

	if result.Error != nil {
		return 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to block user",
			result.Error,
		)
	}

	return result.RowsAffected, nil
}

func (r *repository) Reactivate(ctx context.Context, id string) (int64, error) {
	result := r.DB(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Update("status", security.UserStatusInactive)

	if result.Error != nil {
		return 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to reactivate user",
			result.Error,
		)
	}

	return result.RowsAffected, nil
}

func (r *repository) Activate(ctx context.Context, id string) (int64, error) {
	result := r.DB(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Update("status", security.UserStatusActive)

	if result.Error != nil {
		return 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to activate user",
			result.Error,
		)
	}

	return result.RowsAffected, nil
}
