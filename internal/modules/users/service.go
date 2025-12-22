package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(req CreateUserRequest) (*User, error)
	GetById(id string) (*User, error)
	Delete(id string) error
	Restore(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req CreateUserRequest) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetById(id string) (*User, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.FindById(id)
}

func (s *service) Delete(id string) error {
	return s.repo.SoftDelete(id)
}

func (s *service) Restore(id string) error {
	return s.repo.Restore(id)
}
