package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// SuccessResponse defines the standard structure for successful API responses
type SuccessResponse struct {
	Success bool `json:"success" example:"true"`
	Data    any  `json:"data"`
	Meta    any  `json:"meta,omitempty"`
}

// ErrorResponse defines the standard structure for failed API responses
type ErrorResponse struct {
	Success bool       `json:"success" example:"false"`
	Error   ErrorBlock `json:"error"`
}

// ErrorBlock represents the nested error object details
type ErrorBlock struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty" swaggertype:"object"`
}

// OK sends a 200/201 response wrapped in the SuccessResponse struct
func OK(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// OKWithMeta sends a response including metadata (like pagination)
func OKWithMeta(c *gin.Context, status int, data any, meta any) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Fail sends a manual error response
// Use this when need to construct an error response manually
func Fail(c *gin.Context, status int, code string, message string, fields map[string]string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Error: ErrorBlock{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	})
	c.Abort()
}

// FailWithError maps an internal error to an HTTP error response
// This is the preferred method for handling errors in handlers
// It automatically logs errors and maps them to appropriate HTTP responses
func FailWithError(c *gin.Context, err error) {
	mapped := MapError(err)

	// Log error with request context (request_id is already in context from RequestID middleware)
	logger := log.Ctx(c.Request.Context())

	// Log based on severity:
	// - 5xx errors: log as Error (server issues)
	// - 4xx errors: log as Warn (client issues, but useful for debugging)
	// - Unexpected errors: always log as Error
	logEvent := logger.Error()
	if mapped.Status < 500 {
		logEvent = logger.Warn()
	}

	logEvent.
		Err(err).
		Int("status_code", mapped.Status).
		Str("error_code", mapped.Code).
		Str("path", c.Request.URL.Path).
		Str("method", c.Request.Method).
		Msg("request failed")

	Fail(
		c,
		mapped.Status,
		mapped.Code,
		mapped.Message,
		mapped.Fields,
	)
}
