package internal

import (
	pb "PaymentSystem/protobuf"
	"context"
	"errors"
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
	pb.RegisterUserServiceServer(server, &Server{UserService: service2})
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("failed to start the Server: %v", err)
		}
	}()
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func (m *mockService) CreateUser(ctx context.Context, user CreateUserRequest) (int, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int), args.Error(1)
}

func (m *mockService) GetUser(ctx context.Context, idUser int) (*User, error) {
	args := m.Called(ctx, idUser)
	if args.Get(0).(*User) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *mockService) CreatePixKey(ctx context.Context, pix PixKey) (string, error) {
	args := m.Called(ctx, pix)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}

func (m *mockService) GetPixKey(ctx context.Context, value string) (*GetPixKeyResponse, error) {
	args := m.Called(ctx, value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPixKeyResponse), args.Error(1)
}

func (m *mockService) DeletePixKey(ctx context.Context, value string) error {
	args := m.Called(ctx, value)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return args.Error(0)
}
func (m *mockService) DeleteUser(ctx context.Context, idUser int) error {
	args := m.Called(ctx, idUser)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return args.Error(0)
}

func (m *mockService) UpdateUser(ctx context.Context, user User) (*User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

var _ = Describe("Server Test", func() {
	Context("Server Test", func() {
		var grpcClient pb.UserServiceClient
		var servMock *mockService
		var ctx context.Context
		BeforeEach(func() {
			servMock = new(mockService)
			ctx = context.Background()
			conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(servMock)))
			if err != nil {
				log.Fatal(err)
			}
			grpcClient = pb.NewUserServiceClient(conn)
		})
		It("should CreateUser successfully", func() {
			userPB := &pb.CreateUserRequest{
				Name:  "Name First",
				Age:   20,
				Phone: "+5512912345678",
				Email: "name1@gmail.com",
				Cpf:   "12345678912",
			}
			servMock.On("CreateUser", mock.Anything, CreateUserRequest{
				Name:  "Name First",
				Age:   20,
				Phone: "+5512912345678",
				Email: "name1@gmail.com",
				CPF:   "12345678912",
			}).Return(1, nil)

			res, err := grpcClient.CreateUser(ctx, userPB)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(res.UserId).Should(Equal("1"))
		})
		It("should CreateUser unsuccessfully", func() {
			userPB := &pb.CreateUserRequest{
				Name:  "Name First",
				Age:   20,
				Phone: "+5512912345678",
				Email: "name1@gmail.com",
				Cpf:   "12345678912",
			}
			servMock.On("CreateUser", mock.Anything, CreateUserRequest{
				Name:  "Name First",
				Age:   20,
				Phone: "+5512912345678",
				Email: "name1@gmail.com",
				CPF:   "12345678912",
			}).Return(0, errors.New("error while CreateUser"))

			res, err := grpcClient.CreateUser(ctx, userPB)
			Expect(err).Should(HaveOccurred())
			Expect(res).Should(BeNil())
		})
		It("should GetUser successfully", func() {
			idRequest := &pb.GetUserRequest{UserId: "1"}

			servMock.On("GetUser", mock.Anything, 1).Return(&User{
				Name:    "Name First",
				Age:     20,
				Phone:   "+5512912345678",
				Email:   "name1@gmail.com",
				CPF:     "12345678912",
				Balance: 100,
			}, nil)

			res, err := grpcClient.GetUser(ctx, idRequest)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(res.Age).Should(Equal(int32(20)))
		})
		It("GetUser should fail", func() {
			idRequest := &pb.GetUserRequest{UserId: "1"}

			servMock.On("GetUser", mock.Anything, 1).Return((*User)(nil), errors.New("error while GetUser"))
			res, err := grpcClient.GetUser(ctx, idRequest)
			Expect(err).Should(HaveOccurred())
			Expect(res).Should(BeNil())
		})
		It("should CreatePixKey successfully", func() {
			requiredPix := &pb.CreatePixKeyRequest{
				UserId:   "1",
				KeyType:  "cpf",
				KeyValue: "12345678912",
			}
			servMock.On("CreatePixKey", mock.Anything, mock.MatchedBy(func(req PixKey) bool {
				Expect(req).ShouldNot(BeNil())
				Expect(req.UserID).Should(Equal(1))
				Expect(req.KeyType).Should(Equal(CPF))
				Expect(req.KeyValue).Should(Equal("12345678912"))
				return true
			})).Return("1", nil)
			res, err := grpcClient.CreatePixKey(ctx, requiredPix)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(res.KeyId).Should(Equal("1"))
		})
		It("CreatePixKey should fail", func() {
			requiredPix := &pb.CreatePixKeyRequest{
				UserId:   "1",
				KeyType:  "cpf",
				KeyValue: "12345678912",
			}
			servMock.On("CreatePixKey", mock.Anything, mock.MatchedBy(func(req PixKey) bool {
				Expect(req).ShouldNot(BeNil())
				Expect(req.UserID).Should(Equal(1))
				Expect(req.KeyType).Should(Equal(CPF))
				Expect(req.KeyValue).Should(Equal("12345678912"))
				return true
			})).Return("", errors.New("error while CreatePixKey"))
			res, err := grpcClient.CreatePixKey(ctx, requiredPix)
			Expect(err).Should(HaveOccurred())
			Expect(res).Should(BeNil())
		})
		It("should GetPixKey successfully", func() {
			value := &pb.GetPixKeyRequest{KeyValue: "12345678912"}
			servMock.On("GetPixKey", mock.Anything, "12345678912").Return(&GetPixKeyResponse{
				UserID:   1,
				Name:     "Name Test",
				CPF:      "12345678912",
				KeyID:    "1",
				KeyValue: "123******12",
			}, nil)
			res, err := grpcClient.GetPixKey(context.Background(), value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(res.UserId).Should(Equal("1"))
			Expect(res.KeyId).Should(Equal("1"))
			Expect(res.Name).Should(Equal("Name Test"))
			Expect(res.KeyValue).Should(Equal("123******12"))
			Expect(res.Cpf).Should(Equal("12345678912"))
		})
		It("GetPixKey should fail", func() {
			value := &pb.GetPixKeyRequest{KeyValue: "12345678912"}
			servMock.On("GetPixKey", mock.Anything, "12345678912").Return((*GetPixKeyResponse)(nil), errors.New("error while GetPixKey"))
			res, err := grpcClient.GetPixKey(context.Background(), value)
			Expect(err).Should(HaveOccurred())
			Expect(res).Should(BeNil())
		})
		It("should DeletePixKey successfully", func() {
			request := &pb.DeletePixKeyRequest{KeyValue: "12345678912"}
			servMock.On("DeletePixKey", mock.Anything, request.KeyValue).Return(nil)
			_, err := grpcClient.DeletePixKey(context.Background(), request)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("DeletePixKey should fail", func() {
			request := &pb.DeletePixKeyRequest{KeyValue: "12345678912"}
			servMock.On("DeletePixKey", mock.Anything, request.KeyValue).Return(errors.New("error while DeletePixKey"))
			_, err := grpcClient.DeletePixKey(context.Background(), request)
			Expect(err).Should(HaveOccurred())
		})
		It("should DeleteUser successfully", func() {
			request := &pb.DeleteUserRequest{UserId: "1"}
			servMock.On("DeleteUser", mock.Anything, 1).Return(nil)
			_, err := grpcClient.DeleteUser(context.Background(), request)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("DeleteUser should fail", func() {
			request := &pb.DeleteUserRequest{UserId: "1"}
			servMock.On("DeleteUser", mock.Anything, 1).Return(errors.New("error while DeleteUser"))
			_, err := grpcClient.DeleteUser(context.Background(), request)
			Expect(err).Should(HaveOccurred())
		})
		It(" should UpdateUser successfully", func() {
			request := &pb.UpdateUserRequest{
				UserId: "1",
				Name:   "Name First",
				Age:    20,
				Phone:  "+5512912345678",
				Email:  "name1@gmail.com",
				Cpf:    "12345678912",
			}
			user := User{
				ID:    1,
				Name:  "Name First",
				Age:   20,
				Phone: "+5512912345678",
				Email: "name1@gmail.com",
				CPF:   "12345678912",
			}

			servMock.On("UpdateUser", mock.Anything, user).Return(&user, nil)
			up, err := grpcClient.UpdateUser(context.Background(), request)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(up.UserId).Should(Equal("1"))
		})
		It("UpdateUser should fail", func() {
			request := &pb.UpdateUserRequest{
				UserId: "1",
				Name:   "Name First",
				Age:    20,
				Phone:  "+5512912345678",
				Email:  "name1@gmail.com",
				Cpf:    "12345678912",
			}
			user := User{
				ID:    1,
				Name:  "Name First",
				Age:   20,
				Phone: "+5512912345678",
				Email: "name1@gmail.com",
				CPF:   "12345678912",
			}
			servMock.On("UpdateUser", mock.Anything, user).Return(nil, errors.New("error while UpdateUser"))
			_, err := grpcClient.UpdateUser(context.Background(), request)
			Expect(err).Should(HaveOccurred())
		})
	})
})
