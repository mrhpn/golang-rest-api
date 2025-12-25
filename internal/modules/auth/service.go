package auth

import (
	"context"

	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"github.com/mrhpn/go-rest-api/internal/security"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, email, password string) (*security.TokenPair, *users.User, error)
}

type service struct {
	userService     users.Service
	securityHandler *security.JWTHandler
}

func NewService(userService users.Service, jwtHandler *security.JWTHandler) Service {
	return &service{
		userService:     userService,
		securityHandler: jwtHandler,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (*security.TokenPair, *users.User, error) {
	// 1. get user from User module
	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, ErrInvalidCrendentials
	}

	// 2. verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCrendentials
	}

	// 3. create token pair
	tokens, err := s.securityHandler.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, nil, ErrTokenGeneration
	}

	return tokens, user, nil
}
