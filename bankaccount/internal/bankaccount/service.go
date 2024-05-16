package bankaccount

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	TransferPIX(ctx context.Context, request TransferRequest) error
}

type service struct {
	repBankAccount Repository
}

var validate *validator.Validate

func (s *service) TransferPIX(ctx context.Context, request TransferRequest) error {
	err := validate.Struct(request)
	if err != nil {
		errorMsg := ""
		for _, e := range err.(validator.ValidationErrors) {
			errorMsg += "validation error: field: " + e.Field() + ", value: " + e.Value().(string) + "\n"
		}
		return errors.New(errorMsg)
	}

	return nil
}
