package posts

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/pagination"
)

// Repository defines the persistence operations for post entities.
type postRepository interface {
	Create(ctx context.Context, post *Post) error
	FindByID(ctx context.Context, id string) (*Post, error)
	FindByUserID(ctx context.Context, userID string, opts *pagination.QueryOptions) ([]*Post, int64, error)
	List(ctx context.Context, opts *pagination.QueryOptions) ([]*Post, int64, error)
	Update(ctx context.Context, id string, updates *Post) error
	Delete(ctx context.Context, id string) (int64, error)
}

type service struct {
	repo postRepository
}

// NewService constructs a posts Service with the provided repository.
func NewService(repo postRepository) PostService {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID string, req CreatePostRequest) (*Post, error) {
	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = PostStatusDraft
	}

	// Validate status
	if !IsValidPostStatus(status) {
		return nil, errInvalidStatus
	}

	// Create post
	post := &Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Status:  status,
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, err
	}

	log.Ctx(ctx).Info().
		Str("post_id", post.ID).
		Str("user_id", userID).
		Str("title", req.Title).
		Msg("post created")

	return post, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*Post, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) GetByUserID(ctx context.Context, userID string, opts *pagination.QueryOptions) ([]*Post, *httpx.PaginationMeta, error) {
	posts, total, err := s.repo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, nil, err
	}

	return posts, pagination.BuildMeta(opts, total), nil
}

func (s *service) List(ctx context.Context, opts *pagination.QueryOptions) ([]*Post, *httpx.PaginationMeta, error) {
	posts, total, err := s.repo.List(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	return posts, pagination.BuildMeta(opts, total), nil
}

func (s *service) Update(ctx context.Context, id string, userID string, req UpdatePostRequest) error {
	// Check if post exists and belongs to user
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify ownership
	if post.UserID != userID {
		return errUnauthorized
	}

	// Validate status if provided
	if req.Status != "" && !IsValidPostStatus(req.Status) {
		return errInvalidStatus
	}

	// Prepare updates
	updates := &Post{}
	if req.Title != "" {
		updates.Title = req.Title
	}
	if req.Content != "" {
		updates.Content = req.Content
	}
	if req.Status != "" {
		updates.Status = req.Status
	}

	if err = s.repo.Update(ctx, id, updates); err != nil {
		return err
	}

	log.Ctx(ctx).Info().
		Str("post_id", id).
		Str("user_id", userID).
		Msg("post updated")

	return nil
}

func (s *service) Delete(ctx context.Context, id string, userID string) error {
	// Check if post exists and belongs to user
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify ownership
	if post.UserID != userID {
		return errUnauthorized
	}

	affected, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	if affected == 0 {
		return apperror.ErrNotFound
	}

	log.Ctx(ctx).Info().
		Str("post_id", id).
		Str("user_id", userID).
		Msg("post deleted")

	return nil
}
