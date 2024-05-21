package bankaccount

import (
	pb "PaymentSystem/protobuf"
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedPixServiceServer
	serviceBankAccount Service
}

func (s *Server) Transfer(ctx context.Context, request *pb.TransferRequest) (*emptypb.Empty, error) {
	if request.FromUserId == "" {
		log.Error().Msg("cannot transfer with empty ID")
		return nil, status.Errorf(codes.InvalidArgument, "cannot transfer with empty ID")
	}
	req := ConvertTransferRequest(request)
	err := s.serviceBankAccount.TransferPIX(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("error while transfer pix")
		return nil, status.Errorf(codes.Internal, "error while transfer pix %s", err.Error())
	}
	return &emptypb.Empty{}, nil
}
