package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"github.com/mrhpn/go-rest-api/internal/security"
)

type userProvider interface {
	GetByID(ctx context.Context, id string) (*users.User, error)
	GetByEmail(ctx context.Context, email string) (*users.User, error)
	Activate(ctx context.Context, id string) (*users.User, error)
}

type service struct {
	userProvider    userProvider
	securityHandler *security.JWTHandler
}

// NewService - constructs Auth Service
func NewService(userProvider userProvider, jwtHandler *security.JWTHandler) Service {
	return &service{
		userProvider:    userProvider,
		securityHandler: jwtHandler,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (*security.TokenPair, *users.User, error) {
	// 1. get user from User module
	user, err := s.userProvider.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, errInvalidCrendentials
	}

	// 2. verify password
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, errInvalidCrendentials
	}

	// 3. check if user is blocked
	if user.Status == security.UserStatusBlocked {
		return nil, nil, security.ErrBlockedUser
	}

	// 4. update user status to active on successful login
	if user, err = s.userProvider.Activate(ctx, user.ID); err != nil {
		return nil, nil, err
	}

	// 6. create token pair
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
	user, err := s.userProvider.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", security.ErrInvalidToken
	}

	// 3. check user status. if user is blocked, can't get new token
	if user.Status == security.UserStatusBlocked {
		return "", security.ErrBlockedUser
	}

	// 4. generate ONLY a new access token
	tokens, err := s.securityHandler.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return "", errTokenGeneration
	}

	return tokens.AccessToken, nil
}
