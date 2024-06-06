package bankaccount

import (
	pb "PaymentSystem/protobuf"
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
)

type mockService struct {
	mock.Mock
	Service
}

func dialer(service2 *mockService) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterPixServiceServer(server, &Server{
		serviceBankAccount: service2,
	})
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("failed to start the Server: %v", err)
		}
	}()
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func (m *mockService) TransferPIX(ctx context.Context, request TransferRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *mockService) DepositAmount(ctx context.Context, deposit DepositAmountRequest) error {
	args := m.Called(ctx, deposit)
	return args.Error(0)
}

var _ = Describe("Server Test", func() {
	Context("Transfer", func() {
		var grpcClient pb.PixServiceClient
		var servMock *mockService
		var ctx context.Context
		BeforeEach(func() {
			servMock = new(mockService)
			ctx = context.Background()
			conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(servMock)))
			if err != nil {
				log.Fatal(err)
			}
			grpcClient = pb.NewPixServiceClient(conn)
		})
		It("should Transfer successfully", func() {
			pbRequest := &pb.TransferRequest{
				FromUserId: "1",
				ToPixKey:   "2",
				Amount:     20.00,
			}
			servMock.On("TransferPIX", mock.Anything, mock.MatchedBy(func(req TransferRequest) bool {
				Expect(req).ToNot(BeNil())
				Expect(req.PayerID).Should(Equal("1"))
				Expect(req.ReceiverPixKey).Should(Equal("2"))
				Expect(req.Amount).Should(Equal(20.00))
				return true
			})).Return(nil)
			_, err := grpcClient.Transfer(ctx, pbRequest)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should Transfer unsuccessfully - empty id", func() {
			pbRequest := &pb.TransferRequest{
				FromUserId: "",
				ToPixKey:   "2",
				Amount:     20.00,
			}
			_, err := grpcClient.Transfer(ctx, pbRequest)
			Expect(err).Should(HaveOccurred())
			s, ok := status.FromError(err)
			Expect(s.Code()).To(Equal(codes.InvalidArgument))
			Expect(ok).To(Equal(true))
		})
		It("should Transfer unsuccessfully - empty pix", func() {
			pbRequest := &pb.TransferRequest{
				FromUserId: "1",
				ToPixKey:   "",
				Amount:     20.00,
			}
			_, err := grpcClient.Transfer(ctx, pbRequest)
			Expect(err).Should(HaveOccurred())
			s, ok := status.FromError(err)
			Expect(s.Code()).To(Equal(codes.InvalidArgument))
			Expect(ok).To(Equal(true))
		})
		It("should Transfer unsuccessfully - invalid amount", func() {
			pbRequest := &pb.TransferRequest{
				FromUserId: "1",
				ToPixKey:   "2",
				Amount:     0.00,
			}
			_, err := grpcClient.Transfer(ctx, pbRequest)
			Expect(err).Should(HaveOccurred())
			s, ok := status.FromError(err)
			Expect(s.Code()).To(Equal(codes.InvalidArgument))
			Expect(ok).To(Equal(true))
		})
		It("should Transfer unsuccessfully", func() {
			pbRequest := &pb.TransferRequest{
				FromUserId: "1",
				ToPixKey:   "2",
				Amount:     20.00,
			}
			servMock.On("TransferPIX", mock.Anything, mock.MatchedBy(func(req TransferRequest) bool {
				Expect(req).ToNot(BeNil())
				Expect(req.PayerID).Should(Equal("1"))
				Expect(req.ReceiverPixKey).Should(Equal("2"))
				Expect(req.Amount).Should(Equal(20.00))
				return true
			})).Return(errors.New("error while TransferPIX"))
			_, err := grpcClient.Transfer(ctx, pbRequest)
			Expect(err).Should(HaveOccurred())
		})
	})
	Context("DepositAmount", func() {
		var grpcClient pb.PixServiceClient
		var servMock *mockService
		var ctx context.Context
		BeforeEach(func() {
			servMock = new(mockService)
			ctx = context.Background()
			conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(servMock)))
			if err != nil {
				log.Fatal(err)
			}
			grpcClient = pb.NewPixServiceClient(conn)
		})
		It("should pass DepositAmount successfully", func() {
			pbDeposit := &pb.DepositAmountRequest{
				Amount: "2000.00",
				UserId: "1",
			}
			servMock.On("DepositAmount", mock.Anything, mock.MatchedBy(func(dep DepositAmountRequest) bool {
				Expect(dep).ToNot(BeNil())
				Expect(dep.Amount).Should(Equal(2000.00))
				Expect(dep.UserID).Should(Equal("1"))
				return true
			})).Return(nil)
			_, err := grpcClient.DepositAmount(ctx, pbDeposit)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should pass DepositAmount unsuccessfully - empty id", func() {
			pbDeposit := &pb.DepositAmountRequest{
				Amount: "2000.00",
				UserId: "",
			}
			_, err := grpcClient.DepositAmount(ctx, pbDeposit)
			Expect(err).Should(HaveOccurred())
			s, ok := status.FromError(err)
			Expect(s.Code()).To(Equal(codes.InvalidArgument))
			Expect(ok).To(Equal(true))
		})
		It("should pass DepositAmount unsuccessfully - invalid amount", func() {
			pbDeposit := &pb.DepositAmountRequest{
				Amount: "0.00",
				UserId: "1",
			}
			_, err := grpcClient.DepositAmount(ctx, pbDeposit)
			Expect(err).Should(HaveOccurred())
			s, ok := status.FromError(err)
			Expect(s.Code()).To(Equal(codes.InvalidArgument))
			Expect(ok).To(Equal(true))
		})
		It("should pass DepositAmount unsuccessfully", func() {
			pbDeposit := &pb.DepositAmountRequest{
				Amount: "2000.00",
				UserId: "1",
			}
			servMock.On("DepositAmount", mock.Anything, mock.MatchedBy(func(dep DepositAmountRequest) bool {
				Expect(dep).ToNot(BeNil())
				Expect(dep.UserID).Should(Equal("1"))
				Expect(dep.Amount).Should(Equal(2000.00))
				return true
			})).Return(errors.New("error while DepositAmount"))
			_, err := grpcClient.DepositAmount(ctx, pbDeposit)
			Expect(err).Should(HaveOccurred())
		})
	})
})
