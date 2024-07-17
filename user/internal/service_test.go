package internal

import (
	"PaymentSystem/user/internal/client"
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
	Repository
}

type mockClientBankAccount struct {
	mock.Mock
	client.BankAccount
}

func (m *mockClientBankAccount) GetBalance(ctx context.Context, userID int) (client.Balance, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(client.Balance), args.Error(1)
}

func (m *mockRepository) CreateUser(ctx context.Context, user User) (int, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int), args.Error(1)
}

func (m *mockRepository) GetUser(ctx context.Context, idUser int) (*User, error) {
	args := m.Called(ctx, idUser)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *mockRepository) CreatePixKey(ctx context.Context, pix PixKey) (string, error) {
	args := m.Called(ctx, pix)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}

func (m *mockRepository) GetPixKey(ctx context.Context, value string) (*PixKey, error) {
	args := m.Called(ctx, value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PixKey), args.Error(1)
}

func (m *mockRepository) DeletePixKey(ctx context.Context, idKey string) error {
	args := m.Called(ctx, idKey)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return args.Error(0)
}

func (m *mockRepository) DeleteUser(ctx context.Context, idUser int) error {
	args := m.Called(ctx, idUser)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return args.Error(0)
}

func (m *mockRepository) UpdateUser(ctx context.Context, user User) (*User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

var _ = Describe("Service Test", func() {
	var mockUser *mockRepository
	var mockClient *mockClientBankAccount
	var serv Service
	BeforeEach(func() {
		mockUser = new(mockRepository)
		mockClient = new(mockClientBankAccount)
		serv = NewService(mockUser, mockClient)
	})
	It("should CreateUser successfully", func() {
		user := User{
			Name:    "Name1",
			Age:     20,
			Phone:   "+55 12 91234 5678",
			Email:   "name_1@gmail.com",
			CPF:     "12345678912",
			Balance: 0.0,
		}
		mockUser.On("CreateUser", context.Background(), user).Return(1, nil)
		createUser := CreateUserRequest{
			Name:  "Name1",
			Age:   20,
			Phone: "+55 12 91234 5678",
			Email: "name_1@gmail.com",
			CPF:   "12345678912",
		}
		id, err := serv.CreateUser(context.Background(), createUser)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(id).Should(Equal(1))

	})
	It("should CreateUser unsuccessfully", func() {
		user := User{
			Name:    "",
			Age:     20,
			Phone:   "+55 12 91234 5678",
			Email:   "name_1@gmail.com",
			CPF:     "12345678912",
			Balance: 0.0,
		}
		mockUser.On("CreateUser", context.Background(), user).Return(0, errors.New("error while CreateUser"))
		createUser := CreateUserRequest{
			Name:  "",
			Age:   20,
			Phone: "+55 12 91234 5678",
			Email: "name_1@gmail.com",
			CPF:   "12345678912",
		}
		_, err := serv.CreateUser(context.Background(), createUser)
		Expect(err).Should(HaveOccurred())
	})
	It("should GetUser successfully", func() {
		mockClient.On("GetBalance", context.Background(), 1).Return(client.Balance(2000.00), nil)
		mockUser.On("GetUser", context.Background(), 1).Return(&User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678900",
		}, nil)
		user := &User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678900",
		}
		user, err := serv.GetUser(context.Background(), 1)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(user.ID).Should(Equal(1))
		Expect(user.Age).Should(Equal(int32(30)))
		Expect(user.Balance).Should(Equal(2000.00))
	})
	It("should GetUser unsuccessfully", func() {
		mockClient.On("GetBalance", context.Background(), 1).Return(client.Balance(2000.00), nil)
		mockUser.On("GetUser", context.Background(), 1).Return(nil, errors.New("error while GetUser"))
		user, err := serv.GetUser(context.Background(), 1)
		Expect(err).Should(HaveOccurred())
		Expect(user).To(BeNil())
	})
	It("should GetBalance unsuccessfully", func() {
		mockClient.On("GetBalance", context.Background(), 1).Return(client.Balance(0), errors.New("error while GetBalance"))
		_, err := serv.GetUser(context.Background(), 1)
		Expect(err).Should(HaveOccurred())
	})
	It("should CreatePixKey successfully", func() {
		pix := PixKey{
			UserID:   1,
			KeyType:  CPF,
			KeyValue: "12345678912",
		}
		mockUser.On("CreatePixKey", mock.Anything, pix).Return("1", nil)
		idKey, err := serv.CreatePixKey(context.Background(), pix)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(idKey).Should(Equal("1"))
	})
	It("CreatePixKey should fail", func() {
		pix := PixKey{
			UserID:   1,
			KeyType:  CPF,
			KeyValue: "12345678912",
		}
		mockUser.On("CreatePixKey", mock.Anything, pix).Return("", errors.New("error while CreatePixKey"))
		idKey, err := serv.CreatePixKey(context.Background(), pix)
		Expect(err).Should(HaveOccurred())
		Expect(idKey).Should(Equal(""))
	})
	It("should GetPixKey successfully", func() {
		pix := &PixKey{
			KeyID:    "1",
			UserID:   1,
			KeyType:  CPF,
			KeyValue: "12345678912",
		}
		mockUser.On("GetPixKey", context.Background(), "12345678912").Return(pix, nil)
		mockUser.On("GetUser", context.Background(), pix.UserID).Return(&User{
			ID:    1,
			Name:  "John",
			Email: "johndoe@example.com",
			CPF:   "12345678912",
		}, nil)
		res, err := serv.GetPixKey(context.Background(), "12345678912")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(res.CPF).Should(Equal("123******12"))
		Expect(res.UserID).Should(Equal(1))
		Expect(res.KeyID).Should(Equal("1"))
		Expect(res.Name).Should(Equal("John"))
		Expect(res.KeyValue).Should(Equal("12345678912"))

	})
	It("GetPixKey should fail - while GetUser", func() {
		pix := &PixKey{
			KeyID:    "1",
			UserID:   1,
			KeyType:  CPF,
			KeyValue: "12345678912",
		}
		mockUser.On("GetPixKey", context.Background(), "12345678912").Return(pix, nil)
		mockUser.On("GetUser", context.Background(), 1).Return(nil, errors.New("error while GetUser"))
		res, err := serv.GetPixKey(context.Background(), "12345678912")
		Expect(err).Should(HaveOccurred())
		Expect(res).Should(BeNil())
	})
	It("GetPixKey should fail - while GetPixKey", func() {
		mockUser.On("GetPixKey", context.Background(), "12345678912").Return(nil, errors.New("error while GetPixKey"))
		res, err := serv.GetPixKey(context.Background(), "12345678912")
		Expect(err).Should(HaveOccurred())
		Expect(res).Should(BeNil())
	})
	It("should DeletePixKey successfully", func() {
		pix := &PixKey{
			KeyID:    "1",
			UserID:   1,
			KeyType:  CPF,
			KeyValue: "12345678912",
		}
		mockUser.On("GetPixKey", context.Background(), "12345678912").Return(pix, nil)
		mockUser.On("DeletePixKey", context.Background(), "1").Return(nil)
		err := serv.DeletePixKey(context.Background(), pix.KeyValue)
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("DeletePixKey should fail - while GetPixKey", func() {
		mockUser.On("GetPixKey", context.Background(), "12345678912").Return(nil, errors.New("error while GetPixKey"))
		err := serv.DeletePixKey(context.Background(), "12345678912")
		Expect(err).Should(HaveOccurred())
	})
	It("DeletePixKey should fail - while DeletePixKey", func() {
		pix := &PixKey{
			KeyID:    "1",
			UserID:   1,
			KeyType:  CPF,
			KeyValue: "12345678912",
		}
		mockUser.On("GetPixKey", context.Background(), "12345678912").Return(pix, nil)
		mockUser.On("DeletePixKey", context.Background(), "1").Return(errors.New("while DeletePixKey"))
		err := serv.DeletePixKey(context.Background(), pix.KeyValue)
		Expect(err).Should(HaveOccurred())
	})
	It("should DeleteUser successfully", func() {
		user := &User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678912",
		}
		mockUser.On("GetUser", context.Background(), 1).Return(user, nil)
		mockUser.On("DeleteUser", context.Background(), 1).Return(nil)
		err := serv.DeleteUser(context.Background(), user.ID)
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("DeleteUser should fail - while GetUser", func() {
		user := &User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678912",
		}
		mockUser.On("GetUser", context.Background(), 1).Return(nil, errors.New("error while GetUser"))
		err := serv.DeleteUser(context.Background(), user.ID)
		Expect(err).Should(HaveOccurred())
	})
	It("DeleteUser should fail - while DeleteUser", func() {
		user := &User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678912",
		}
		mockUser.On("GetUser", context.Background(), 1).Return(user, nil)
		mockUser.On("DeleteUser", context.Background(), 1).Return(errors.New("error while DeleteUser"))
		err := serv.DeleteUser(context.Background(), user.ID)
		Expect(err).Should(HaveOccurred())
	})
	It("should UpdateUser successfully", func() {
		user := User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678912",
		}
		upUser := &User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678912",
		}
		mockUser.On("GetUser", context.Background(), user.ID).Return(upUser, nil)
		mockUser.On("UpdateUser", context.Background(), user).Return(upUser, nil)
		up, err := serv.UpdateUser(context.Background(), user)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(up).Should(Equal(upUser))
	})
	It("UpdateUser should fail - while GetUser", func() {
		user := User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678912",
		}
		mockUser.On("GetUser", context.Background(), 1).Return(nil, errors.New("error while GetUser"))
		up, err := serv.UpdateUser(context.Background(), user)
		Expect(err).Should(HaveOccurred())
		Expect(up).Should(BeNil())
	})
	It("UpdateUser should fail - while UpdateUser", func() {
		user := User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678912",
		}
		upUser := &User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678912",
		}
		mockUser.On("GetUser", context.Background(), user.ID).Return(upUser, nil)
		mockUser.On("UpdateUser", context.Background(), user).Return(nil, errors.New("error while UpdateUser"))
		up, err := serv.UpdateUser(context.Background(), user)
		Expect(err).Should(HaveOccurred())
		Expect(up).Should(BeNil())
	})
})
