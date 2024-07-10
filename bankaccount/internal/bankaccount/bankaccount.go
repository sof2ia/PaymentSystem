package bankaccount

import (
	pb "PaymentSystem/protobuf"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strconv"
)

type TransferRequest struct {
	PayerID        string  `validate:"required"`
	ReceiverPixKey string  `validate:"required"`
	Amount         float64 `validate:"required,gt=0,lte=5000"`
}

var validate *validator.Validate

func (t *TransferRequest) Validation() error {
	validate = validator.New()
	err := validate.Struct(t)
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
	return nil
}

type DepositAmountRequest struct {
	Amount float64 `validate:"required,gt=0,lte=5000"`
	UserID string  `validate:"required"`
}

func (d *DepositAmountRequest) Validation() error {
	validate = validator.New()
	err := validate.Struct(d)
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
	return nil
}

type TransferResponse struct {
	Success bool
}

type Movement struct {
	ID            string
	Amount        string
	UserID        string
	Date          string
	TransactionID string
	OperationType OperationType
}

type OperationType string

const (
	Credit OperationType = "credit"
	Debt   OperationType = "debt"
)

func ConvertTransferRequest(requestPB *pb.TransferRequest) (TransferRequest, error) {
	t := TransferRequest{
		PayerID:        requestPB.FromUserId,
		ReceiverPixKey: requestPB.ToPixKey,
		Amount:         requestPB.Amount,
	}
	err := t.Validation()
	if err != nil {
		return TransferRequest{}, err
	}
	return t, nil
}

func ConvertDepositAmount(depositPB *pb.DepositAmountRequest) (DepositAmountRequest, error) {
	depositAmount, err := strconv.ParseFloat(depositPB.Amount, 64)
	if err != nil {
		return DepositAmountRequest{}, err
	}
	d := DepositAmountRequest{
		Amount: depositAmount,
		UserID: depositPB.UserId,
	}
	err = d.Validation()
	if err != nil {
		return DepositAmountRequest{}, err
	}
	return d, nil
}

type GetBalanceRequest struct {
	ID int
}

type GetBalanceResponse struct {
	Balance float64
}
