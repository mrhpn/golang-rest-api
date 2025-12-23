package httpx

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	appErr "github.com/mrhpn/go-rest-api/internal/errors"
	"github.com/mrhpn/go-rest-api/internal/utils"
)

func BindJSON(c *gin.Context, req any) error {
	if err := c.ShouldBindJSON(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fields := make(map[string]string)
			for _, fe := range ve {
				fields[utils.ToSnakeCase(fe.Field())] = utils.GetValidationMessage(fe)
			}

			return &appErr.AppError{
				Kind:    appErr.InvalidInput,
				Code:    "INVALID_REQUEST",
				Message: "invalid request",
				Fields:  fields,
			}
		}

		return appErr.New(
			appErr.BadRequest,
			"BAD_REQUEST",
			"invalid request",
		)
	}
	return nil
}

func BindURI(c *gin.Context, req any) error {
	if err := c.ShouldBindUri(req); err != nil {
		return appErr.New(
			appErr.InvalidInput,
			"INVALID_URI",
			"invalid recource id",
		)
	}
	return nil
}
