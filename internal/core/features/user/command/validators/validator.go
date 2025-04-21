package validators

import (
	"github.com/go-playground/validator/v10"
)

// TODO: is it better to move it to utils ?
func ValidateRequest(req interface{}) error {
	validate := validator.New()
	return validate.Struct(req)
}
