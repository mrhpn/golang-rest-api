package httpx

import "github.com/gin-gonic/gin"

type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
	Meta    any  `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Success bool       `json:"success"`
	Error   ErrorBlock `json:"error"`
}

type ErrorBlock struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func OKWithMeta(c *gin.Context, status int, data any, meta any) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func Fail(c *gin.Context, status int, code string, message string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Error: ErrorBlock{
			Code:    code,
			Message: message,
		},
	})
}
