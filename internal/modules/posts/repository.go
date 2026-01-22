package posts

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/pagination"
	repo "github.com/mrhpn/go-rest-api/internal/repository"
)

type Repository struct {
	repo.Base
}

// NewRepository constructs a posts Repository backed by a GORM database.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Base: repo.Base{
			DBInstance: db,
		},
	}
}

func (r *Repository) Create(ctx context.Context, post *Post) error {
	err := r.DB(ctx).Create(post).Error
	if err != nil {
		return apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to create post",
			err,
		)
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*Post, error) {
	var post Post
	err := r.DB(ctx).
		Preload("User").
		First(&post, "id = ?", id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrNotFound
		}
		return nil, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to find post",
			err,
		)
	}

	return &post, nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID string, opts *pagination.QueryOptions) ([]*Post, int64, error) {
	var posts []*Post
	var total int64

	// 1. Get total count
	err := r.DB(ctx).Model(&Post{}).
		Where("user_id = ?", userID).
		Scopes(pagination.SearchScope(opts)).
		Count(&total).Error
	if err != nil {
		return nil, 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to count posts",
			err,
		)
	}

	// 2. Fetch data
	err = r.DB(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Scopes(pagination.Paginate(opts)).
		Find(&posts).Error
	if err != nil {
		return nil, 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to find posts",
			err,
		)
	}

	return posts, total, nil
}

func (r *Repository) List(ctx context.Context, opts *pagination.QueryOptions) ([]*Post, int64, error) {
	var posts []*Post
	var total int64

	// 1. Get total count
	err := r.DB(ctx).Model(&Post{}).
		Scopes(pagination.SearchScope(opts)).
		Count(&total).Error
	if err != nil {
		return nil, 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to count posts",
			err,
		)
	}

	// 2. Fetch data
	err = r.DB(ctx).
		Preload("User").
		Scopes(pagination.Paginate(opts)).
		Find(&posts).Error
	if err != nil {
		return nil, 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to find posts",
			err,
		)
	}

	return posts, total, nil
}

func (r *Repository) Update(ctx context.Context, id string, updates *Post) error {
	result := r.DB(ctx).
		Model(&Post{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to update post",
			result.Error,
		)
	}

	if result.RowsAffected == 0 {
		return apperror.ErrNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) (int64, error) {
	result := r.DB(ctx).Delete(&Post{}, "id = ?", id)
	if result.Error != nil {
		return 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to delete post",
			result.Error,
		)
	}
	return result.RowsAffected, nil
}
