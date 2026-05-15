package validator

import (
	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate = validator.New()

func Validates(i any) error {
	err := Validator.Struct(i)

	if err != nil {
		//sends in the complete map of fields
		errors := err.(validator.ValidationErrors)
		return errors
	}
	return nil
}
