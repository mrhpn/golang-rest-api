package users

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, req CreateUserRequest) (*User, error)
	GetById(ctx context.Context, id string) (*User, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	// check email uniqueness
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailExists
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, err // unexpected DB error
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, ErrInternal
	}

	// create user
	user := &User{
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	log.Ctx(ctx).Info().Str("email", req.Email).Msg("user created")

	return user, nil
}

func (s *service) GetById(ctx context.Context, id string) (*User, error) {
	return s.repo.FindById(ctx, id)
}

func (s *service) Delete(ctx context.Context, id string) error {
	affected, err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrUserNotFound
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
		return ErrUserNotFound
	}

	log.Ctx(ctx).Info().Str("user_id", id).Msg("user restored")

	return nil
}
