package internal

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	CreateUser(ctx context.Context, user CreateUserRequest) (int64, error)
}

type service struct {
	userRepository Repository
}

func (s *service) CreateUser(ctx context.Context, user CreateUserRequest) (userID int64, err error) {
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

func NewService(userRepository Repository) Service {
	return &service{userRepository: userRepository}
}
