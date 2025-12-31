package httpx

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

// RegisterValidators registers validators used in http-request steps
func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validators here
		_ = v.RegisterValidation("ulid", validateULID)
	}
}

func validateULID(fl validator.FieldLevel) bool {
	ulidStr := fl.Field().String()
	_, err := ulid.Parse(ulidStr)
	return err == nil
}
