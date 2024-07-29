package client

import (
	"context"
	pb "github.com/sof2ia/PaymentSystem/bankaccount/protobuf"
	"strconv"
)

type BankAccount interface {
	GetBalance(ctx context.Context, userID int) (Balance, error)
}

type Balance float64

type bankAccountClient struct {
	client pb.PixServiceClient
}

func (b *bankAccountClient) GetBalance(ctx context.Context, userID int) (Balance, error) {
	userIDStr := strconv.Itoa(userID)
	res, err := b.client.GetBalance(ctx, &pb.GetBalanceRequest{IdUser: userIDStr})
	if err != nil {
		return 0, err
	}
	return Balance(res.Balance), nil
}

func NewBankAccountClient(client pb.PixServiceClient) BankAccount {
	return &bankAccountClient{client: client}
}
