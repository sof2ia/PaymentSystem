package internal

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	CreateUser(ctx context.Context, user User) error
}

type service struct {
	userRepository Repository
}

func (s *service) CreateUser(ctx context.Context, user User) error {
	err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return status.Errorf(codes.Internal, "error while CreateUser %s", err)
	}
	return nil
}

func NewService(userRepository Repository) Service {
	return &service{userRepository: userRepository}
}
