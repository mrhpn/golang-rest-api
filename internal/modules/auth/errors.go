package auth

import "github.com/mrhpn/go-rest-api/internal/apperror"

var (
	// ErrIdentityNotFoundInContext indicates that the authenticated identity is missing from the request context.
	ErrIdentityNotFoundInContext = apperror.New(
		apperror.Internal,
		"INTERNAL_ERROR",
		"identity not found in context",
	)

	// ErrUnauthorized indicates that the request lacks valid authentication credentials, such as a missing access token.
	ErrUnauthorized = apperror.New(
		apperror.Unauthorized,
		"UNAUTHORIZED",
		"unauthorized: missing token",
	)

	// ErrForbidden indicates that the authenticated identity does not
	// have sufficient permissions to access the requested resource.
	ErrForbidden = apperror.New(
		apperror.Forbidden,
		"FORBIDDEN",
		"forbidden: insufficient permissions",
	)

	// ErrRefreshTokenMissing indicates that a refresh token was expected but not provided in the request.
	errRefreshTokenMissing = apperror.New(
		apperror.Unauthorized,
		"REFRESH_TOKEN_MISSING",
		"refresh token missing",
	)

	// ErrInvalidCrendentials indicates that the provided authentication credentials are invalid.
	errInvalidCrendentials = apperror.New(
		apperror.Unauthorized,
		"INVALID_CRENDENTIALS",
		"invalid crendentials",
	)

	// ErrTokenGeneration indicates a failure during access or refresh token generation.
	errTokenGeneration = apperror.New(
		apperror.Internal,
		"ERR_TOKEN_GENERATION",
		"cannot generate token",
	)
)
