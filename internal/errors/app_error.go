package errors

import "fmt"

type Kind string

const (
	NotFound     Kind = "NOT_FOUND"
	Unauthorized Kind = "UNAUTHORIZED"
	Forbidden    Kind = "FORBIDDEN"
	Conflict     Kind = "CONFLICT"
	InvalidInput Kind = "INVALID_INPUT" // prioritized over BadRequest
	BadRequest   Kind = "BAD_REQUEST"
	Internal     Kind = "INTERNAL"
)

// AppError represents an application-level error with structured information
// for proper HTTP response mapping and client communication
type AppError struct {
	Kind    Kind              // for HTTP mapping
	Code    string            // frontend error code
	Message string            // human-readable message (safe for clients)
	Fields  map[string]string // optional field-specific errors (used in validations)
	Err     error             // underlying error (for logging, not exposed to clients)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error for error wrapping support
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError without an underlying error
func New(kind Kind, code string, message string) *AppError {
	return &AppError{
		Kind:    kind,
		Code:    code,
		Message: message,
	}
}

// Wrap creates a new AppError wrapping an underlying error
// The underlying error is used for logging but not exposed to clients
func Wrap(kind Kind, code string, message string, err error) *AppError {
	return &AppError{
		Kind:    kind,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WithFields adds field-specific validation errors to an AppError
func (e *AppError) WithFields(fields map[string]string) *AppError {
	e.Fields = fields
	return e
}
