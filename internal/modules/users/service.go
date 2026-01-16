package users

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/pagination"
	"github.com/mrhpn/go-rest-api/internal/security"
)

// Service defines the business logic for managing users.
type Service interface {
	Create(ctx context.Context, req CreateUserRequest) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, opts *pagination.QueryOptions) ([]*User, *httpx.PaginationMeta, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
	Block(ctx context.Context, id string) error
	Reactivate(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

// NewService constructs a users Service with the provided repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	// check email uniqueness
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, errEmailExists
	}
	if !errors.Is(err, errUserNotFound) {
		return nil, err // unexpected DB error
	}

	// hash password
	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.Wrap(
			apperror.Internal,
			apperror.ErrInternal.Code,
			"failed to hash password",
			err,
		)
	}

	// create user
	user := &User{
		Email:        req.Email,
		Role:         req.Role,
		PasswordHash: hash,
	}

	if user.Role == "" {
		user.Role = security.RoleEmployee
	}

	if err = s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	log.Ctx(ctx).Info().
		Str("email", req.Email).
		Str("role", string(user.Role)).
		Msg("user created")

	return user, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *service) List(ctx context.Context, opts *pagination.QueryOptions) ([]*User, *httpx.PaginationMeta, error) {
	users, total, err := s.repo.List(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	return users, pagination.BuildMeta(opts, total), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	affected, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	if affected == 0 {
		return errUserNotFound
	}

	log.Ctx(ctx).Info().Str("user_id", id).Msg("user deleted")

	return nil
}

func (s *service) Restore(ctx context.Context, id string) error {
	affected, err := s.repo.Restore(ctx, id)
	if err != nil {
		return err
	}

	if affected == 0 {
		return errUserNotFound
	}

	log.Ctx(ctx).Info().Str("user_id", id).Msg("user restored")

	return nil
}

func (s *service) Block(ctx context.Context, id string) error {
	affected, err := s.repo.Block(ctx, id)
	if err != nil {
		return err
	}

	if affected == 0 {
		return errUserNotFound
	}

	log.Ctx(ctx).Info().Str("user_id", id).Msg("user blocked")

	return nil
}

func (s *service) Reactivate(ctx context.Context, id string) error {
	affected, err := s.repo.Reactivate(ctx, id)
	if err != nil {
		return err
	}

	if affected == 0 {
		return errUserNotFound
	}

	log.Ctx(ctx).Info().Str("user_id", id).Msg("user reactivated")

	return nil
}

func (s *service) Activate(ctx context.Context, id string) error {
	affected, err := s.repo.Activate(ctx, id)
	if err != nil {
		return err
	}

	if affected == 0 {
		return errUserNotFound
	}

	log.Ctx(ctx).Info().Str("user_id", id).Msg("user activated")
	return nil
}
