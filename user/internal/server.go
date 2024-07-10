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
	userPB := &pb.CreateUserResponse{UserId: strconv.Itoa(userID)}
	return userPB, nil
}

func (s *Server) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	idUser, err := strconv.Atoi(request.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	user, err := s.userService.GetUser(ctx, idUser)
	if err != nil {
		log.Error().Err(err).Msg("error while GetUser")
		return nil, status.Errorf(codes.Internal, "error while GetUser %s", err.Error())
	}
	userPB, err := ConvertGetUserResponse(user)
	if err != nil {
		log.Error().Err(err).Msg("error while ConvertGetUserResponse")
		return nil, status.Errorf(codes.Internal, "error while ConvertGetUserResponse %s", err.Error())
	}
	return userPB, nil
}

func (s *Server) CreatePixKey(ctx context.Context, request *pb.CreatePixKeyRequest) (*pb.CreatePixKeyResponse, error) {
	idUser, err := strconv.Atoi(request.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	requiredPix := PixKey{
		UserID:   idUser,
		KeyType:  KeyType(request.KeyType),
		KeyValue: request.KeyValue,
	}
	newPix, err := s.userService.CreatePixKey(ctx, requiredPix)
	if err != nil {
		log.Error().Err(err).Msg("error while CreatePixKey")
		return nil, status.Errorf(codes.Internal, "error while CreatePixKey %s", err.Error())
	}
	pixPB := &pb.CreatePixKeyResponse{KeyId: newPix}
	return pixPB, nil
}
