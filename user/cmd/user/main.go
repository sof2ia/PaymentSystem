package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	pb "github.com/sof2ia/PaymentSystem/bankaccount/protobuf"
	"github.com/sof2ia/PaymentSystem/user/internal"
	"github.com/sof2ia/PaymentSystem/user/internal/client"
	pb2 "github.com/sof2ia/PaymentSystem/user/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
)

func main() {
	conn, err := startPostgreSQL(context.Background())
	if err != nil {
		return
	}
	defer conn.Close(context.Background())
	// Create the gRPC client
	grpcConn, err := grpc.DialContext(context.Background(), os.Getenv("GRPC_SERVICES_SVC"),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start connection with grpc")
	}
	defer func(grpcConn *grpc.ClientConn) {
		err = grpcConn.Close()
		if err != nil {
			log.Err(err).Msg("Failed to close grpc connection")
		}
	}(grpcConn)

	pixClient := pb.NewPixServiceClient(grpcConn)
	bankAccountClient := client.NewBankAccountClient(pixClient)
	userRep := internal.NewRepository(conn)
	userServ := internal.NewService(userRep, bankAccountClient)
	userSer := internal.Server{UserService: userServ}

	grpcServer := grpc.NewServer()
	pb2.RegisterUserServiceServer(grpcServer, &userSer)
	list, err := net.Listen("tcp", ":9001")
	log.Printf("Start on port 9001")
	err = grpcServer.Serve(list)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to serve gRPC server on port")
	}
}

func startPostgreSQL(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database:")
		return nil, err
	}
	return conn, nil
}
