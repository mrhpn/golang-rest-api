package security

import "github.com/mrhpn/go-rest-api/internal/errors"

var (
	ErrInvalidToken = errors.New(
		errors.BadRequest,
		"INVALID_TOKEN",
		"invalid token",
	)

	ErrExpiredToken = errors.New(
		errors.BadRequest,
		"EXPIRED_TOKEN",
		"token has expired",
	)

	ErrBlockedUser = errors.New(
		errors.Unauthorized,
		"USER_BLOCKED",
		"user is blocked",
	)
)
