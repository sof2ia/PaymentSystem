package internal

import (
	"PaymentSystem/user/internal/client"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	CreateUser(ctx context.Context, user CreateUserRequest) (int, error)
	GetUser(cxt context.Context, idUser int) (*User, error)
	CreatePixKey(ctx context.Context, pix PixKey) (string, error)
}

type service struct {
	userRepository    Repository
	bankAccountClient client.BankAccount
}

func (s *service) CreatePixKey(ctx context.Context, pix PixKey) (string, error) {
	idKey, err := s.userRepository.CreatePixKey(ctx, pix)
	if err != nil {
		return "", status.Errorf(codes.Internal, "error while CreatePixKey %s, %s", pix.KeyType, pix.KeyValue)
	}
	return idKey, err
}

func (s *service) GetUser(ctx context.Context, idUser int) (*User, error) {
	balance, err := s.bankAccountClient.GetBalance(ctx, idUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while GetBalance %v", balance)
	}
	user, err := s.userRepository.GetUser(ctx, idUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while GetUser %s", err)
	}
	user.Balance = float64(balance)
	return user, nil
}

func (s *service) CreateUser(ctx context.Context, user CreateUserRequest) (userID int, err error) {
	createUser := User{
		Name:    user.Name,
		Age:     user.Age,
		Phone:   user.Phone,
		Email:   user.Email,
		CPF:     user.CPF,
		Balance: 0.0,
	}
	userID, err = s.userRepository.CreateUser(ctx, createUser)
	if err != nil {
		return 0, status.Errorf(codes.Internal, "error while CreateUser %s", err)
	}
	return userID, nil
}

func NewService(userRepository Repository, bankAccountClient client.BankAccount) Service {
	return &service{userRepository: userRepository, bankAccountClient: bankAccountClient}
}
