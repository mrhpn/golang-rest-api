package posts

import "github.com/mrhpn/go-rest-api/internal/apperror"

var (
	// errUserNotFound indicates that a requested user does not exist.
	errInvalidStatus = apperror.New(
		apperror.BadRequest,
		"INVALID_STATUS",
		"invalid post status",
	)

	// errUnauthorized indicates that a requested user can't modify the resource
	errUnauthorized = apperror.New(
		apperror.Unauthorized,
		"UNAUTHORIZED",
		"unauthorized to modify this resource",
	)
)
