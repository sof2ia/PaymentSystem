package internal

import (
	pb "PaymentSystem/protobuf"
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
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
	pb.RegisterPixServiceServer(server, &Server{userService: service2})
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("failed to start the Server: %v", err)
		}
	}()
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func (m *mockService) CreateUser(ctx context.Context, user CreateUserRequest) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

var _ = Describe("Server Test", func() {
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
	It("should CreateUser successfully", func() {
		userPB := &pb.CreateUserRequest{
			Name:  "Name First",
			Age:   20,
			Phone: "+5512912345678",
			Email: "name1@gmail.com",
			Cpf:   "12345678912",
		}
		servMock.On("CreateUser", context.Background(), CreateUserRequest{
			Name:  "Name First",
			Age:   20,
			Phone: "+5512912345678",
			Email: "name1@gmail.com",
			CPF:   "12345678912",
		}).Return(int64(1), nil)

		res, err := grpcClient.CreateUser(ctx, userPB)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(res.AccountId).Should(Equal("1"))
	})
})
