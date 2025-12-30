package auth

import "github.com/mrhpn/go-rest-api/internal/errors"

var (
	ErrInternal = errors.New(
		errors.Internal,
		"INTERNAL_ERROR",
		"internal server error",
	)

	ErrInvalidID = errors.New(
		errors.InvalidInput,
		"INVALID_ID_FORMAT",
		"invalid id format",
	)

	ErrIdentityNotFoundInContext = errors.New(
		errors.Internal,
		"INTERNAL_ERROR",
		"identity not found in context",
	)

	ErrUnauthorized = errors.New(
		errors.Unauthorized,
		"UNAUTHORIZED",
		"unauthorized: missing token",
	)

	ErrForbidden = errors.New(
		errors.Forbidden,
		"FORBIDDEN",
		"forbidden: insufficient permissions",
	)

	ErrRefreshTokenMissing = errors.New(
		errors.Unauthorized,
		"REFRESH_TOKEN_MISSING",
		"refresh token missing",
	)

	ErrInvalidCrendentials = errors.New(
		errors.Unauthorized,
		"INVALID_CRENDENTIALS",
		"invalid crendentials",
	)

	ErrTokenGeneration = errors.New(
		errors.Internal,
		"ERR_TOKEN_GENERATION",
		"cannot generate token",
	)
)
