package errors

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

type AppError struct {
	Kind    Kind              // for HTTP mapping
	Code    string            // frontend error code
	Message string            // human-readable message
	Fields  map[string]string // optional field-specific errors (used in validations)
}

func (e *AppError) Error() string {
	return e.Message
}

func New(kind Kind, code string, message string) *AppError {
	return &AppError{
		Kind:    kind,
		Code:    code,
		Message: message,
	}
}
