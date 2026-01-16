// Package apperror provides utilities for working with application errors.
package apperror

import "fmt"

// Kind represents the category of an application error.
type Kind string

const (
	// NotFound indicates that a requested resource does not exist.
	NotFound Kind = "NOT_FOUND"
	// Unauthorized indicates missing or invalid authentication.
	Unauthorized Kind = "UNAUTHORIZED"
	// Forbidden indicates insufficient permissions.
	Forbidden Kind = "FORBIDDEN"
	// Conflict indicates a state conflict, such as a duplicate resource.
	Conflict Kind = "CONFLICT"
	// InvalidInput indicates validation failures or invalid input based on business logic.
	InvalidInput Kind = "INVALID_INPUT"
	// BadRequest indicates a malformed or invalid request.
	BadRequest Kind = "BAD_REQUEST"
	// Internal indicates an unexpected internal server error.
	Internal Kind = "INTERNAL"
	// TooManyRequests indicates client sends too many requests to server.
	TooManyRequests Kind = "RATE_LIMIT_EXCEEDED"
	// RequestTimeout indicates client's request is timeout and execution will be stopped.
	RequestTimeout Kind = "REQUEST_TIMEOUT"
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

// Error returns a human-readable string representation of the error.
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
