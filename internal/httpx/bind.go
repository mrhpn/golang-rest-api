package httpx

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/stringx"
	"github.com/mrhpn/go-rest-api/internal/validation"
)

// BindJSON binds the request body to the given struct.
func BindJSON(c *gin.Context, req any) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return handleBindingError(err)
	}
	return nil
}

// BindURI binds the request URI to the given struct.
func BindURI(c *gin.Context, req any) error {
	if err := c.ShouldBindUri(req); err != nil {
		return apperror.New(
			apperror.InvalidInput,
			"INVALID_URI",
			"invalid recource id",
		)
	}
	return nil
}

// BindQuery binds the query parameters to the given struct.
func BindQuery(c *gin.Context, req any) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return handleBindingError(err)
	}
	return nil
}

func handleBindingError(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make(map[string]string)
		for _, fe := range ve {
			fields[stringx.ToSnakeCase(fe.Field())] = validation.GetValidationMessage(fe)
		}

		return &apperror.AppError{
			Kind:    apperror.InvalidInput,
			Code:    "INVALID_REQUEST",
			Message: "invalid request",
			Fields:  fields,
		}
	}

	return apperror.New(
		apperror.BadRequest,
		"BAD_REQUEST",
		"failed to parse request",
	)
}
