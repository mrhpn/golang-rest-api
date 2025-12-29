package httpx

import "github.com/gin-gonic/gin"

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
func Fail(c *gin.Context, status int, code string, message string, fields map[string]string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Error: ErrorBlock{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	})
}

// FailWithError maps an internal error to an HTTP error response
func FailWithError(c *gin.Context, err error) {
	mapped := MapError(err)

	Fail(
		c,
		mapped.Status,
		mapped.Code,
		mapped.Message,
		mapped.Fields,
	)
}
