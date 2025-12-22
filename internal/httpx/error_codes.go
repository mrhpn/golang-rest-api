package httpx

const (
	// General errors
	ErrBadRequest   = "BAD_REQUEST"
	ErrNotFound     = "NOT_FOUND"
	ErrInternal     = "INTERNAL_ERROR"
	ErrUnauthorized = "UNAUTHORIZED"
	ErrForbidden    = "FORBIDDEN"

	// User module errors
	ErrUserNotFound = "USER_NOT_FOUND"
	ErrEmailExists  = "EMAIL_ALREADY_EXISTS"
)
