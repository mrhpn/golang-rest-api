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
