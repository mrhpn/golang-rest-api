package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/mrhpn/go-rest-api/internal/stringx"
)

type messageFormatter func(field string, fe validator.FieldError) string

//nolint:gochecknoglobals // immutable validation message registry
var validationMessages = map[string]messageFormatter{
	// ---------- presence ----------
	"required": func(field string, _ validator.FieldError) string {
		return fmt.Sprintf("%s is required", field)
	},

	// ---------- format ----------
	"email": func(_ string, _ validator.FieldError) string {
		return "invalid email format"
	},

	"uuid": func(_ string, _ validator.FieldError) string {
		return "invalid id format"
	},

	"ulid": func(_ string, _ validator.FieldError) string {
		return "invalid id format"
	},

	"url": func(_ string, _ validator.FieldError) string {
		return "invalid url format"
	},

	"ip": func(_ string, _ validator.FieldError) string {
		return "invalid ip address"
	},

	// ---------- size ----------
	"min": func(field string, fe validator.FieldError) string {
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	},

	"max": func(field string, fe validator.FieldError) string {
		return fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
	},

	"len": func(field string, fe validator.FieldError) string {
		return fmt.Sprintf("%s must be exactly %s characters", field, fe.Param())
	},

	// ---------- numbers ----------
	"gt": func(field string, fe validator.FieldError) string {
		return fmt.Sprintf("%s must be greater than %s", field, fe.Param())
	},

	"gte": func(field string, fe validator.FieldError) string {
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
	},

	"lt": func(field string, fe validator.FieldError) string {
		return fmt.Sprintf("%s must be less than %s", field, fe.Param())
	},

	"lte": func(field string, fe validator.FieldError) string {
		return fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
	},

	// ---------- collections ----------
	"oneof": func(field string, fe validator.FieldError) string {
		return fmt.Sprintf("%s must be one of [%s]", field, fe.Param())
	},
}

// GetValidationMessage returns a human-readable message for a validation error.
func GetValidationMessage(fe validator.FieldError) string {
	field := stringx.ToSnakeCase(fe.Field())

	if formatter, ok := validationMessages[fe.Tag()]; ok {
		return formatter(field, fe)
	}

	return "invalid value"
}
