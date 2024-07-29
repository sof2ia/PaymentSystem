package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rs/zerolog/log"
	"github.com/sof2ia/PaymentSystem/bankaccount/internal"
	pb "github.com/sof2ia/PaymentSystem/bankaccount/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
)

func main() {
	conn, err := startDynamoDB()
	if err != nil {
		return
	}
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

	repBA := internal.NewRepository(conn)
	servBA := internal.NewService(repBA)
	serBA := internal.Server{ServiceBankAccount: servBA}

	grpcServer := grpc.NewServer()
	pb.RegisterPixServiceServer(grpcServer, &serBA)
	list, err := net.Listen("tcp", ":9002")
	log.Printf("Start on port 9002")
	err = grpcServer.Serve(list)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to serve gRPC server on port")
	}
}

func startDynamoDB() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Error().Err(err).Msg("Failed to load default config")
		return nil, err
	}
	clientDynamodb := dynamodb.NewFromConfig(cfg)
	return clientDynamodb, nil
}
