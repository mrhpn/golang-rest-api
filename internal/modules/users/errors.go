package users

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

	ErrUserNotFound = errors.New(
		errors.NotFound,
		"USER_NOT_FOUND",
		"user not found",
	)

	ErrEmailExists = errors.New(
		errors.Conflict,
		"EMAIL_EXISTS",
		"email already exists",
	)

	ErrInvalidUserInput = errors.New(
		errors.InvalidInput,
		"INVALID_USER_INPUT",
		"invalid user input",
	)
)
