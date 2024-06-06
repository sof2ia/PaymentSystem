package internal

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/rs/zerolog/log"
)

var _ = Describe("CreateUser", func() {
	var mock pgxmock.PgxPoolIface
	var rep Repository
	var err error
	BeforeEach(func() {
		mock, err = pgxmock.NewPool()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		defer mock.Close()
		log.Printf("mock was created successfully")
		rep = NewRepository(mock)
	})

	It("should pass successfully", func() {
		user := User{
			Name:    "John Doe",
			Age:     30,
			Phone:   "123456789",
			Email:   "johndoe@example.com",
			CPF:     "12345678900",
			Balance: 0.0,
		}
		mock.ExpectQuery("INSERT INTO Users").WithArgs(user.Name, user.Age, user.Phone, user.Email, user.CPF, user.Balance).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
		id, err := rep.CreateUser(context.Background(), user)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(id).Should(Equal(int64(1)))

	})
	It("should pass unsuccessfully", func() {
		user := User{
			Name:    "",
			Age:     30,
			Phone:   "123456789",
			Email:   "johndoe@example.com",
			CPF:     "12345678900",
			Balance: 0.0,
		}
		mock.ExpectQuery("INSERT INTO Users").WithArgs(user.Name, user.Age, user.Phone, user.Email, user.CPF, user.Balance).WillReturnError(errors.New("fields should not be empty"))

		_, err := rep.CreateUser(context.Background(), user)
		Expect(err).Should(HaveOccurred())
	})
})
