package internal

import (
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

func (m *mockRepository) CreateUser(ctx context.Context, user User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

var _ = Describe("Service Test", func() {
	var mock *mockRepository
	var serv Service
	BeforeEach(func() {
		mock = new(mockRepository)
		serv = NewService(mock)
	})
	It("should CreateUser successfully", func() {
		user := User{
			Name:    "Name1",
			Age:     20,
			Phone:   "+55 12 91234 5678",
			Email:   "name_1@gmail.com",
			CPF:     "123.456.789-12",
			Balance: 0.0,
		}
		mock.On("CreateUser", context.Background(), user).Return(int64(1), nil)
		createUser := CreateUserRequest{
			Name:  "Name1",
			Age:   20,
			Phone: "+55 12 91234 5678",
			Email: "name_1@gmail.com",
			CPF:   "123.456.789-12",
		}
		id, err := serv.CreateUser(context.Background(), createUser)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(id).Should(Equal(int64(1)))

	})
	It("should CreateUser unsuccessfully", func() {
		user := User{
			Name:    "",
			Age:     20,
			Phone:   "+55 12 91234 5678",
			Email:   "name_1@gmail.com",
			CPF:     "123.456.789-12",
			Balance: 0.0,
		}
		mock.On("CreateUser", context.Background(), user).Return(int64(0), errors.New("error while CreateUser"))
		createUser := CreateUserRequest{
			Name:  "",
			Age:   20,
			Phone: "+55 12 91234 5678",
			Email: "name_1@gmail.com",
			CPF:   "123.456.789-12",
		}
		_, err := serv.CreateUser(context.Background(), createUser)
		Expect(err).Should(HaveOccurred())
	})
})
