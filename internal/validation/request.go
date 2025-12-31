package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mrhpn/go-rest-api/internal/stringx"
)

// GetValidationMessage returns a human-readable message for a validation error.
func GetValidationMessage(fe validator.FieldError) string {
	field := stringx.ToSnakeCase(fe.Field())

	switch fe.Tag() {
	// ---------- presence ----------
	case "required":
		return fmt.Sprintf("%s is required", field)

	// ---------- format ----------
	case "email":
		return "invalid email format"

	case "uuid", "ulid":
		return "invalid id format"

	case "url":
		return "invalid url format"

	case "ip":
		return "invalid ip address"

	// ---------- size ----------
	case "min":
		return fmt.Sprintf(
			"%s must be at least %s characters",
			field, fe.Param(),
		)

	case "max":
		return fmt.Sprintf(
			"%s must be at most %s characters",
			field, fe.Param(),
		)

	case "len":
		return fmt.Sprintf(
			"%s must be exactly %s characters",
			field, fe.Param(),
		)

	// ---------- numbers ----------
	case "gt":
		return fmt.Sprintf(
			"%s must be greater than %s",
			field, fe.Param(),
		)

	case "gte":
		return fmt.Sprintf(
			"%s must be greater than or equal to %s",
			field, fe.Param(),
		)

	case "lt":
		return fmt.Sprintf(
			"%s must be less than %s",
			field, fe.Param(),
		)

	case "lte":
		return fmt.Sprintf(
			"%s must be less than or equal to %s",
			field, fe.Param(),
		)

	// ---------- collections ----------
	case "oneof":
		return fmt.Sprintf(
			"%s must be one of [%s]",
			field, fe.Param(),
		)

	// ---------- fallback ----------
	default:
		return "invalid value"
	}
}
