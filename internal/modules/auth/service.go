package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, email, password string) (string, *users.User, error)
}

type service struct {
	userService users.Service
	jwtSecret   string
}

func NewService(userService users.Service, jwtSecret string) Service {
	return &service{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (string, *users.User, error) {
	// get user from User module
	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, ErrInvalidCrendentials
	}

	// verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCrendentials
	}

	// create jwt claims
	claims := middlewares.UserClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// sign token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, ErrTokenGeneration
	}

	return tokenString, user, nil
}
