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

func (m *mockRepository) CreateUser(ctx context.Context, user User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
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
			Name:  "Name1",
			Age:   20,
			Phone: "+55 12 91234 5678",
			Email: "name_1@gmail.com",
			CPF:   "123.456.789-12",
		}
		mock.On("CreateUser", context.Background(), user).Return(nil)
		err := serv.CreateUser(context.Background(), user)
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("should CreateUser unsuccessfully", func() {
		user := User{
			Name:  "",
			Age:   20,
			Phone: "+55 12 91234 5678",
			Email: "name_1@gmail.com",
			CPF:   "123.456.789-12",
		}
		mock.On("CreateUser", context.Background(), user).Return(errors.New("error while CreateUser"))
		err := serv.CreateUser(context.Background(), user)
		Expect(err).Should(HaveOccurred())
	})

})
