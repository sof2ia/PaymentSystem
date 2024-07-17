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
	pb.UnimplementedPixServiceServer
	ServiceBankAccount Service
}

func (s *Server) Transfer(ctx context.Context, request *pb.TransferRequest) (*emptypb.Empty, error) {
	log.Info().Msgf("Starting Transfer: %v", request)
	req, err := ConvertTransferRequest(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	err = s.ServiceBankAccount.TransferPIX(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("error while transfer pix")
		return nil, status.Errorf(codes.Internal, "error while transfer pix %s", err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DepositAmount(ctx context.Context, deposit *pb.DepositAmountRequest) (*emptypb.Empty, error) {
	dep, err := ConvertDepositAmount(deposit)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	err = s.ServiceBankAccount.DepositAmount(ctx, dep)
	if err != nil {
		log.Error().Err(err).Msg("error while DepositAmount")
		return nil, status.Errorf(codes.InvalidArgument, "error while DepositAmount %s", err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetBalance(ctx context.Context, request *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	idUserInt, err := strconv.Atoi(request.IdUser)
	log.Printf("error 1: %d", idUserInt)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	balance, err := s.ServiceBankAccount.GetBalance(ctx, GetBalanceRequest{ID: idUserInt})
	log.Printf("error 2: %v", balance)
	if err != nil {
		log.Error().Err(err).Msg("error while GetBalance")
		return nil, status.Errorf(codes.InvalidArgument, "error while GetBalance %s", err.Error())
	}
	balancePB := &pb.GetBalanceResponse{Balance: balance.Balance}
	log.Printf("error 3: %v", balancePB)
	return balancePB, nil
}
