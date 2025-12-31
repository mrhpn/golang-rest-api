// Package httpx provides utilities for binding request data to structs and validating input.
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
			"invalid request",
		)
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
