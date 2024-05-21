package bankaccount

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/ksuid"
	"strconv"
	"time"
)

type Service interface {
	TransferPIX(ctx context.Context, request TransferRequest) error
}

type service struct {
	repBankAccount Repository
}

var validate *validator.Validate

func (s *service) TransferPIX(ctx context.Context, request TransferRequest) error {

	validate = validator.New()

	err := validate.Struct(request)
	if err != nil {
		errorMsg := ""
		for _, e := range err.(validator.ValidationErrors) {
			if e.Field() == "Amount" {
				errorMsg += "validation error: field: " + e.Field() + ", value: " + fmt.Sprintf("%f", e.Value().(float64)) + "\n"
			} else {
				errorMsg += "validation error: field: " + e.Field() + ", value: " + e.Value().(string) + "\n"
			}
		}
		return errors.New(errorMsg)
	}
	idTransaction := ksuid.New().String()
	debtAmount := strconv.FormatFloat(-request.Amount, 'f', 2, 64)
	creditAmount := strconv.FormatFloat(request.Amount, 'f', 2, 64)

	payerMovement := Movement{
		ID:            ksuid.New().String(),
		Amount:        debtAmount,
		UserID:        request.PayerID,
		Date:          time.Now().Format(time.RFC3339),
		TransactionID: idTransaction,
		OperationType: Debit,
	}
	err = s.repBankAccount.Save(ctx, payerMovement)
	if err != nil {
		return err
	}

	receiverMovement := Movement{
		ID:            ksuid.New().String(),
		Amount:        creditAmount,
		UserID:        request.ReceiverPixKey,
		Date:          time.Now().Format(time.RFC3339),
		TransactionID: idTransaction,
		OperationType: Credit,
	}
	err = s.repBankAccount.Save(ctx, receiverMovement)
	if err != nil {
		return err
	}
	return nil
}

func NewService(repBankAccount Repository) Service {
	return &service{repBankAccount: repBankAccount}
}
