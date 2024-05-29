package internal

import (
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/rs/zerolog/log"
)

var _ = Describe("CreateUser", func() {
	var mock pgxmock.PgxPoolIface
	BeforeEach(func() {
		mock, err := pgxmock.NewPool()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		defer mock.Close()
		rep := NewRepository(mock)
	})
	//mockDB, mock, err := sqlmock.New()
	//if err != nil {
	//	t.Fatalf("the creation of mock is failed %v", err)
	//}
	//defer func(mockDB *sql.DB) {
	//	err := mockDB.Close()
	//	if err != nil {
	//
	//	}
	//}(mockDB)

	It("should pass successfully", func() {
		mock.ExpectExec("INSERT INTO Users").WithArgs("Name1", 20, "+55 12 91234 5678", "name_1@gmail.com", "123.456.789-12").WillReturnResult(sqlmock.NewResult(1, 1))
		result := sqlmock.NewRows([]string{"AccountID", "Name", "Age", "Phone", "Email", "CPF", "Balance"}).AddRow("1", "Name1", 20, "+55 12 91234 5678", "name_1@gmail.com", "123.456.789-12", 0.0)
		mock.ExpectQuery("SELECT \\* FROM Users WHERE ID = ?").WithArgs(1).WillReturnRows(result)
		_, err = rep.CreateUser(User{
			Name:  "Name1",
			Age:   20,
			Phone: "+55 12 91234 5678",
			Email: "name_1@gmail.com",
			CPF:   "123.456.789-12",
		})
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("should pass unsuccessfully", func() {
		mock.ExpectExec("INSERT INTO Users").WithArgs("Name1", 20, "+55 12 91234 5678", "name_1@gmail.com", "123.456.789-12").WillReturnResult(sqlmock.NewResult(1, 1))
		result := sqlmock.NewRows([]string{"AccountID", "Name", "Age", "Phone", "Email", "CPF", "Balance"}).AddRow("1", "Name1", 20, "+55 12 91234 5678", "name_1@gmail.com", "123.456.789-12", 0.0)
		mock.ExpectQuery("SELECT \\* FROM Users WHERE ID = ?").WithArgs(1).WillReturnRows(result)
		_, err = rep.CreateUser(User{
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
