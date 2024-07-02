package internal

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/rs/zerolog/log"
	"regexp"
)

var _ = Describe("Repository Test", func() {
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
	It(" CreateUser should pass successfully", func() {
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
	It("CreateUser should pass unsuccessfully", func() {
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
	It("GetUser should pass successfully", func() {
		user := User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678900",
		}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM Users WHERE ID = $1`)).WithArgs(user.ID).WillReturnRows(pgxmock.NewRows([]string{"id", "name", "age", "phone", "email", "cpf"}).AddRow(1, "John", 30, "123456789", "johndoe@example.com", "12345678900"))

		user, err = rep.GetUser(context.Background(), int64(1))
		Expect(err).ShouldNot(HaveOccurred())
		Expect(user.ID).Should(Equal(int64(1)))
		Expect(user.Age).Should(Equal(30))
	})
})
