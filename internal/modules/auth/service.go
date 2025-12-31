package auth

import (
	"context"

	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"github.com/mrhpn/go-rest-api/internal/security"
	"golang.org/x/crypto/bcrypt"
)

// Service - auth service interface
type Service interface {
	Login(ctx context.Context, email, password string) (*security.TokenPair, *users.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

type service struct {
	userService     users.Service
	securityHandler *security.JWTHandler
}

// NewService - constructs Auth Service
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
		return nil, nil, errInvalidCrendentials
	}

	// 2. verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, errInvalidCrendentials
	}

	// 3. create token pair
	tokens, err := s.securityHandler.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, nil, errTokenGeneration
	}

	return tokens, user, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// 1. validate the signature and expiry of the refresh token
	claims, err := s.securityHandler.ValidateToken(refreshToken)
	if err != nil {
		return "", security.ErrExpiredToken
	}

	// 2. hit the db: ensure the user still exists and isn't blocked
	user, err := s.userService.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", security.ErrInvalidToken
	}

	// todo:
	// if user.status == "Blocked" { return "", security.ErrBlockedUser }

	// 3. generate ONLY a new access token
	tokens, err := s.securityHandler.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return "", errTokenGeneration
	}

	return tokens.AccessToken, nil
}
