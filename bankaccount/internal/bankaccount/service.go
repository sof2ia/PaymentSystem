package bankaccount

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
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
	payerMovements, err := s.repBankAccount.ListMovementsByUser(ctx, request.PayerID)
	if err != nil {
		return err
	}

	var balance float64
	for _, movement := range payerMovements {
		balanceStr, err := strconv.ParseFloat(movement.Amount, 64)
		if err != nil {
			return err
		}
		balance += balanceStr
	}

	if balance < request.Amount {
		return errors.New("insufficient balance")
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
		OperationType: Debt,
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

func (s *service) DepositAmount(ctx context.Context, deposit DepositAmountRequest) error {

	payerMovements, err := s.repBankAccount.ListMovementsByUser(ctx, deposit.UserID)
	if err != nil {
		return err
	}

	var balance float64
	for _, movement := range payerMovements {
		balanceStr, err := strconv.ParseFloat(movement.Amount, 64)
		if err != nil {
			return err
		}
		balance += balanceStr
	}

	newBalance := balance + deposit.Amount
	depositAmount := strconv.FormatFloat(newBalance, 'f', 2, 64)

	depositMovement := Movement{
		ID:            ksuid.New().String(),
		Amount:        depositAmount,
		UserID:        deposit.UserID,
		Date:          time.Now().Format(time.RFC3339),
		TransactionID: ksuid.New().String(),
		OperationType: Credit,
	}
	log.Info().Msgf("amount: %s", depositMovement.Amount)
	err = s.repBankAccount.Save(ctx, depositMovement)
	if err != nil {
		return err
	}
	return nil
}

func NewService(repBankAccount Repository) Service {
	return &service{repBankAccount: repBankAccount}
}
