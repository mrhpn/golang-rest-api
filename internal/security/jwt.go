package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwt payload
type UserClaims struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTHandler struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTHandler(secret string, accessTokenExpirySecond, refreshTokenExpirySecond int) *JWTHandler {
	return &JWTHandler{
		secret:        []byte(secret),
		accessExpiry:  time.Duration(accessTokenExpirySecond) * time.Second,
		refreshExpiry: time.Duration(refreshTokenExpirySecond) * time.Second,
	}
}

// generate access & refresh tokens
func (h *JWTHandler) GenerateTokenPair(userID string, role Role) (*TokenPair, error) {
	accessToken, err := h.signToken(userID, role, h.accessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := h.signToken(userID, role, h.refreshExpiry)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// private: sign token
func (h *JWTHandler) signToken(userID string, role Role, expiry time.Duration) (string, error) {
	claims := UserClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.secret)
}

// validate token
func (h *JWTHandler) ValidateToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return h.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
