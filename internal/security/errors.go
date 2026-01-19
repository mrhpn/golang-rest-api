package security

import "github.com/mrhpn/go-rest-api/internal/apperror"

var (
	// ErrInvalidToken indicates that the provided authentication token is malformed, unverifiable, or otherwise invalid.
	ErrInvalidToken = apperror.New(
		apperror.Unauthorized,
		"INVALID_TOKEN",
		"invalid token",
	)

	// ErrExpiredToken indicates that the provided authentication token has expired and is no longer valid.
	ErrExpiredToken = apperror.New(
		apperror.Unauthorized,
		"EXPIRED_TOKEN",
		"token has expired",
	)

	// ErrBlockedUser indicates that the authenticated user is blocked and is not allowed to access protected resources.
	ErrBlockedUser = apperror.New(
		apperror.Unauthorized,
		"USER_BLOCKED",
		"user is blocked",
	)

	// ErrRequestTimeout indicates that the requests coming from client is too many and server blocked for a period of time.
	ErrRequestTimeout = apperror.New(
		apperror.RequestTimeout,
		"REQUEST_TIMEOUT",
		"request timeout exceeded",
	)

	// ErrForbiddenInProd indicates that this resource is only available in development environment.
	ErrForbiddenInProd = apperror.New(
		apperror.Forbidden,
		"FORBIDDEN_IN_PROD",
		"this operation is only allowed in development env",
	)

	// ErrRouteNotFound indicates that the requested route does not exist or is unavailable.
	ErrRouteNotFound = apperror.New(
		apperror.NotFound,
		"ROUTE_NOT_FOUND",
		"the requested route does not existed or is unavailable",
	)

	// ErrMethodNotAllowed indicates that the requested method does not exist for the requested route.
	ErrMethodNotAllowed = apperror.New(
		apperror.BadRequest,
		"METHOD_NOT_ALLOWED",
		"method not allowed for this route",
	)
)
