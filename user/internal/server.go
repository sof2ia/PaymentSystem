package internal

import (
	pb "PaymentSystem/protobuf"
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	UserService Service
}

func (s *Server) CreateUser(ctx context.Context, user *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	req, err := ConvertCreateUserRequest(user)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	userID, err := s.UserService.CreateUser(ctx, req)
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
	user, err := s.UserService.GetUser(ctx, idUser)
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

func (s *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	up, err := ConvertUpdateUserRequest(request)
	user, err := s.UserService.UpdateUser(ctx, up)
	if err != nil {
		log.Error().Err(err).Msg("error while update user")
		return nil, status.Errorf(codes.Internal, "error while update user %s", err.Error())
	}
	userPB, err := ConvertUpdateUserResponse(user)
	if err != nil {
		log.Error().Err(err).Msg("error while ConvertUpdateUserResponse")
		return nil, status.Errorf(codes.Internal, "error while ConvertUpdateUserResponse %s", err.Error())
	}
	return userPB, nil
}

func (s *Server) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	idUserInt, err := strconv.Atoi(request.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	err = s.UserService.DeleteUser(ctx, idUserInt)
	if err != nil {
		log.Error().Err(err).Msg("error while delete user")
		return nil, status.Errorf(codes.Internal, "error while delete user %s", err.Error())
	}
	return &emptypb.Empty{}, nil
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
	newPix, err := s.UserService.CreatePixKey(ctx, requiredPix)
	if err != nil {
		log.Error().Err(err).Msg("error while CreatePixKey")
		return nil, status.Errorf(codes.Internal, "error while CreatePixKey %s", err.Error())
	}
	pixPB := &pb.CreatePixKeyResponse{KeyId: newPix}
	return pixPB, nil
}

func (s *Server) GetPixKey(ctx context.Context, request *pb.GetPixKeyRequest) (*pb.GetPixKeyResponse, error) {
	pix, err := s.UserService.GetPixKey(ctx, request.KeyValue)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	resPB := &pb.GetPixKeyResponse{
		UserId:   strconv.Itoa(pix.UserID),
		Name:     pix.Name,
		Cpf:      pix.CPF,
		KeyId:    pix.KeyID,
		KeyValue: pix.KeyValue,
	}
	return resPB, nil
}

func (s *Server) DeletePixKey(ctx context.Context, request *pb.DeletePixKeyRequest) (*emptypb.Empty, error) {
	err := s.UserService.DeletePixKey(ctx, request.KeyValue)
	if err != nil {
		log.Error().Err(err).Msg("error while delete pix key")
		return nil, status.Errorf(codes.Internal, "error while delete pix key %s", err.Error())
	}
	return &emptypb.Empty{}, nil
}
