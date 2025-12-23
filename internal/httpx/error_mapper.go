package httpx

import (
	"errors"
	"net/http"

	domain "github.com/mrhpn/go-rest-api/internal/errors"
	"gorm.io/gorm"
)

type mappedError struct {
	Status  int               // http status code
	Code    string            // frontend error code
	Message string            // human-readable message
	Fields  map[string]string // optional field-specific errors (used in validations)
}

func MapError(err error) mappedError {
	var appErr *domain.AppError

	// 1. Check for custom application errors
	if errors.As(err, &appErr) {
		return mappedError{
			Status:  mapKindToStatus(appErr.Kind),
			Code:    appErr.Code,
			Message: appErr.Message,
			Fields:  appErr.Fields,
		}
	}

	// 2. Fallback: catch statndard GORM errors that escaped the repo
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return mappedError{
			Status:  http.StatusNotFound,
			Code:    "NOT_FOUND",
			Message: "resource not found",
		}
	}

	// 3. Ultimate Fallback: actual 500
	return mappedError{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_ERROR",
		Message: "❌ internal server error",
	}
}

func mapKindToStatus(kind domain.Kind) int {
	switch kind {
	case domain.NotFound:
		return http.StatusNotFound
	case domain.InvalidInput, domain.BadRequest:
		return http.StatusBadRequest
	case domain.Conflict:
		return http.StatusConflict
	case domain.Unauthorized:
		return http.StatusUnauthorized
	case domain.Forbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
