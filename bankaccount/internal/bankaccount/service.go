package bankaccount

import (
	"context"
	"github.com/segmentio/ksuid"
	"strconv"
	"time"
)

type Service interface {
	TransferPIX(ctx context.Context, request TransferRequest) error
	DepositAmount(ctx context.Context, deposit DepositAmountRequest) error
}

type service struct {
	repBankAccount Repository
}

func (s *service) TransferPIX(ctx context.Context, request TransferRequest) error {

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
	err := s.repBankAccount.Save(ctx, payerMovement)
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

func (s *service) DepositAmount(ctx context.Context, deposit DepositAmountRequest) error {

	depositAmount := strconv.FormatFloat(deposit.Amount, 'f', 2, 64)
	depositMovement := Movement{
		ID:            ksuid.New().String(),
		Amount:        depositAmount,
		UserID:        deposit.UserID,
		Date:          time.Now().Format(time.RFC3339),
		TransactionID: ksuid.New().String(),
		OperationType: Credit,
	}
	err := s.repBankAccount.Save(ctx, depositMovement)
	if err != nil {
		return err
	}
	return nil
}

func NewService(repBankAccount Repository) Service {
	return &service{repBankAccount: repBankAccount}
}
