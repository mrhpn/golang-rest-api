package users

import "github.com/mrhpn/go-rest-api/internal/apperror"

var (
	// errUserNotFound indicates that a requested user does not exist.
	errUserNotFound = apperror.New(
		apperror.NotFound,
		"USER_NOT_FOUND",
		"user not found",
	)

	// errEmailExists indicates that the provided email address is already associated with an existing user.
	errEmailExists = apperror.New(
		apperror.Conflict,
		"EMAIL_EXISTS",
		"email already exists",
	)
)
