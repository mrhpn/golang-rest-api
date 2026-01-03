package httpx

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/mrhpn/go-rest-api/internal/apperror"
)

type mappedError struct {
	Status  int               // http status code
	Code    string            // frontend error code
	Message string            // human-readable message
	Fields  map[string]string // optional field-specific errors (used in validations)
}

// MapError maps an error to an HTTP response structure
// It sanitizes internal errors and ensures client-safe messages
// Note: Error logging is handled in FailWithError to include request context
func mapError(err error) mappedError {
	var appErr *apperror.AppError

	// 1. Check for custom application errors
	if errors.As(err, &appErr) {
		return mappedError{
			Status:  mapKindToStatus(appErr.Kind),
			Code:    appErr.Code,
			Message: appErr.Message,
			Fields:  appErr.Fields,
		}
	}

	// 2. Fallback: catch standard GORM errors that escaped the repo
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return mappedError{
			Status:  http.StatusNotFound,
			Code:    "NOT_FOUND",
			Message: "resource not found",
		}
	}

	// 3. Ultimate Fallback: actual 500
	// This should rarely happen if all errors are properly wrapped
	return mappedError{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
	}
}

func mapKindToStatus(kind apperror.Kind) int {
	switch kind {
	case apperror.NotFound:
		return http.StatusNotFound
	case apperror.InvalidInput, apperror.BadRequest:
		return http.StatusBadRequest
	case apperror.Conflict:
		return http.StatusConflict
	case apperror.Unauthorized:
		return http.StatusUnauthorized
	case apperror.Forbidden:
		return http.StatusForbidden
	case apperror.TooManyRequests:
		return http.StatusTooManyRequests
	case apperror.Internal:
		return http.StatusInternalServerError
	}
	return http.StatusInternalServerError
}
