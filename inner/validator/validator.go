package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type RequestValidator struct {
	validate *validator.Validate
}

func NewRequestValidator() *RequestValidator {
	validate := validator.New()
	return &RequestValidator{validate: validate}
}

func (v RequestValidator) Validate(request any) (err error) {
	err = v.validate.Struct(request)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return validateErrs
		}
	}
	return err
}
