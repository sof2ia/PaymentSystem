package internal

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var validate *validator.Validate

func (c *CreateUserRequest) Validation() error {
	validate = validator.New()
	err := validate.RegisterValidation("CPF", validateCPF)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("phone", validatePhone)
	if err != nil {
		return err
	}
	err = validate.Struct(c)
	if err != nil {
		errorMsg := ""
		for _, e := range err.(validator.ValidationErrors) {
			if e.Field() == "Age" {
				errorMsg += "validation error: field: " + e.Field() + ", value: " + fmt.Sprintf("%d", e.Value().(int32)) + "\n"
			} else {
				errorMsg += "validation error: field: " + e.Field() + ", value: " + e.Value().(string) + "\n"
			}
		}
		return errors.New(errorMsg)
	}
	return nil
}

func validateCPF(fl validator.FieldLevel) bool {
	isValid, err := regexp.MatchString(`^\d{11}$`, fl.Field().String())
	if err != nil {
		return false
	}
	return isValid
}

func validatePhone(fl validator.FieldLevel) bool {
	isValid, err := regexp.MatchString(`^\+55\d{2}9\d{8}$`, fl.Field().String())
	if err != nil {
		return false
	}
	return isValid
}
