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
	UpdateUser(ctx context.Context, user User) (*User, error)
	DeleteUser(ctx context.Context, idUser int) error
	CreatePixKey(ctx context.Context, pix PixKey) (string, error)
	GetPixKey(ctx context.Context, value string) (*GetPixKeyResponse, error)
	DeletePixKey(ctx context.Context, value string) error
}

type service struct {
	userRepository    Repository
	bankAccountClient client.BankAccount
}

func (s *service) UpdateUser(ctx context.Context, user User) (*User, error) {
	_, err := s.userRepository.GetUser(ctx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user was not found")
	}

	updateUser, err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while UpdateUser %s", err)
	}

	return updateUser, nil
}

func (s *service) DeleteUser(ctx context.Context, idUser int) error {
	_, err := s.userRepository.GetUser(ctx, idUser)
	if err != nil {
		return status.Errorf(codes.Internal, "user was not found")
	}

	err = s.userRepository.DeleteUser(ctx, idUser)
	if err != nil {
		return status.Errorf(codes.Internal, "error while DeleteUser %s", err)
	}

	return nil
}

func (s *service) DeletePixKey(ctx context.Context, value string) error {
	key, err := s.userRepository.GetPixKey(ctx, value)
	if err != nil {
		return status.Errorf(codes.Internal, "pix key was not found")
	}

	err = s.userRepository.DeletePixKey(ctx, key.KeyID)
	if err != nil {
		return status.Errorf(codes.Internal, "error while DeletePixKey %s", key.KeyID)
	}

	return nil
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
func (s *service) GetPixKey(ctx context.Context, value string) (*GetPixKeyResponse, error) {
	pix, err := s.userRepository.GetPixKey(ctx, value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while GetPixKey, %s", err)
	}
	user, err := s.userRepository.GetUser(ctx, pix.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while GetUser, %s", err)
	}
	response := &GetPixKeyResponse{
		UserID:   user.ID,
		Name:     user.Name,
		CPF:      maskCPF(user.CPF),
		KeyID:    pix.KeyID,
		KeyValue: pix.KeyValue,
	}
	return response, nil
}

func maskCPF(cpf string) string {
	mask := cpf[:3] + "******" + cpf[9:]
	return mask
}

func NewService(userRepository Repository, bankAccountClient client.BankAccount) Service {
	return &service{userRepository: userRepository, bankAccountClient: bankAccountClient}
}
