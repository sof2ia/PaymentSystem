package internal

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/rs/zerolog/log"
)

var _ = Describe("CreateUser", func() {
	var mock pgxmock.PgxPoolIface
	var rep Repository
	BeforeEach(func() {
		mock, err := pgxmock.NewPool()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		log.Printf("mockn was created successfully")
		defer mock.Close()
		rep = NewRepository(mock)
	})
	It("should pass successfully", func() {
		mock.ExpectExec("INSERT INTO Users").WithArgs("Name1", 20, "+55 12 91234 5678", "name_1@gmail.com", "123.456.789-12").WillReturnResult(pgxmock.NewResult("INSERT", 1))
		err := rep.CreateUser(context.Background(), User{
			Name:  "Name1",
			Age:   20,
			Phone: "+55 12 91234 5678",
			Email: "name_1@gmail.com",
			CPF:   "123.456.789-12",
		})
		Expect(err).ShouldNot(HaveOccurred())

	})
	It("should pass unsuccessfully", func() {
		mock.ExpectExec("INSERT INTO Users").WithArgs("Name1", 20, "+55 12 91234 5678", "name_1@gmail.com", "123.456.789-12").WillReturnResult(pgxmock.NewResult("INSERT", 1))
		err := rep.CreateUser(context.Background(), User{
			Name:  "",
			Age:   20,
			Phone: "+55 12 91234 5678",
			Email: "name_1@gmail.com",
			CPF:   "123.456.789-12",
		})
		log.Err(err)
		Expect(err).Should(HaveOccurred())
	})
})
