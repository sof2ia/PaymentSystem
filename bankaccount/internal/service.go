package internal

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
	GetBalance(ctx context.Context, request GetBalanceRequest) (*GetBalanceResponse, error)
}

type service struct {
	repBankAccount Repository
}

func (s *service) GetBalance(ctx context.Context, request GetBalanceRequest) (*GetBalanceResponse, error) {
	requestIDStr := strconv.Itoa(request.ID)
	log.Printf("error 4: %s", requestIDStr)
	payerMovements, err := s.repBankAccount.ListMovementsByUser(ctx, requestIDStr)
	log.Printf("error 5: %+v", payerMovements)
	if err != nil {
		return nil, err
	}

	var balance float64
	for _, movement := range payerMovements {
		balanceStr, err := strconv.ParseFloat(movement.Amount, 64)
		if err != nil {
			return nil, err
		}
		balance += balanceStr
	}
	log.Printf("error 6: %v", balance)

	return &GetBalanceResponse{Balance: balance}, nil
}

func (s *service) TransferPIX(ctx context.Context, request TransferRequest) error {
	payerIDStr, err := strconv.Atoi(request.PayerID)
	if err != nil {
		return err
	}
	balance, err := s.GetBalance(ctx, GetBalanceRequest{ID: payerIDStr})
	if err != nil {
		return err
	}
	if balance.Balance < request.Amount {
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

	depositIDStr, err := strconv.Atoi(deposit.UserID)
	if err != nil {
		return err
	}
	balance, err := s.GetBalance(ctx, GetBalanceRequest{ID: depositIDStr})
	if err != nil {
		return err
	}

	newBalance := balance.Balance + deposit.Amount
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
