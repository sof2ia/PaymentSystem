package internal

import (
	pb "PaymentSystem/protobuf"
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type Server struct {
	pb.UnimplementedPixServiceServer
	userService Service
}

func (s *Server) CreateUser(ctx context.Context, user *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	req, err := ConvertCreateUserRequest(user)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	userID, err := s.userService.CreateUser(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("error while CreateUser")
		return nil, status.Errorf(codes.Internal, "error while CreateUser %s", err.Error())
	}
	userPB := &pb.CreateUserResponse{AccountId: strconv.FormatInt(userID, 10)}
	return userPB, nil
}
