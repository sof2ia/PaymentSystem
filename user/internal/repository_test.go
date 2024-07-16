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
		Expect(id).Should(Equal(1))

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
		user := &User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678900",
		}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM Users WHERE ID = $1`)).WithArgs(user.ID).WillReturnRows(pgxmock.NewRows([]string{"id", "name", "age", "email", "phone", "cpf"}).AddRow(user.ID, user.Name, user.Age, user.Email, user.Phone, user.CPF))

		newUser, err := rep.GetUser(context.Background(), 1)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(newUser).Should(Equal(user))
		Expect(newUser.ID).Should(Equal(1))
		Expect(newUser.Age).Should(Equal(int32(30)))
	})
	It("GetUser should fail", func() {
		user := &User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678900",
		}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM Users WHERE ID = $1`)).WithArgs(user.ID).WillReturnError(errors.New("error while GetUser"))
		newUser, err := rep.GetUser(context.Background(), 1)
		Expect(err).Should(HaveOccurred())
		Expect(newUser).ShouldNot(Equal(user))
	})
	It("CreatePixKey should pass successfully", func() {
		pix := PixKey{
			UserID:   1,
			KeyType:  CPF,
			KeyValue: "123.456.789-12",
		}
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO PixKey`)).WithArgs(pix.UserID, pix.KeyType, pix.KeyValue).WillReturnRows(pgxmock.NewRows([]string{"key_id"}).AddRow("1"))
		newPix, err := rep.CreatePixKey(context.Background(), pix)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(newPix).Should(Equal("1"))
	})
	It("CreatePixKey should fail", func() {
		pix := PixKey{
			UserID:   1,
			KeyType:  CPF,
			KeyValue: "123.456.789-12",
		}
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO PixKey`)).WithArgs(pix.UserID, pix.KeyType, pix.KeyValue).WillReturnError(errors.New("error while CreatePixKey"))
		_, err := rep.CreatePixKey(context.Background(), pix)
		Expect(err).Should(HaveOccurred())
	})
	It("should GetPixKey pass successfully", func() {
		resPix := &PixKey{
			KeyID:    "1",
			UserID:   1,
			KeyType:  CPF,
			KeyValue: "123.456.789-12",
		}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM PixKey WHERE key_value = $1`)).WithArgs("123.456.789-12").WillReturnRows(pgxmock.NewRows([]string{"key_id", "user_id", "key_type", "key_value"}).AddRow(resPix.KeyID, resPix.UserID, resPix.KeyType, resPix.KeyValue))
		pix, err := rep.GetPixKey(context.Background(), "123.456.789-12")
		Expect(pix).Should(Equal(resPix))
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("GetPixKey should fail", func() {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM PixKey WHERE key_value = $1`)).WithArgs("123.456.789-12").WillReturnError(errors.New("error while GetPixKey"))
		pix, err := rep.GetPixKey(context.Background(), "123.456.789-12")
		Expect(err).Should(HaveOccurred())
		Expect(pix).Should(BeNil())
	})
	It("should UpdateUser successfully", func() {
		user := User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678900",
		}
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE Users SET Name = $1, Age = $2, Phone = $3, Email = $4, CPF = $5 WHERE ID =$6`)).WithArgs(user.Name, user.Age, user.Phone, user.Email, user.CPF, user.ID).WillReturnResult(pgxmock.NewResult("", 1))
		up, err := rep.UpdateUser(context.Background(), user)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(up.ID).Should(Equal(1))
		Expect(up.Name).Should(Equal("John"))
		Expect(up.Age).Should(Equal(int32(30)))
		Expect(up.Phone).Should(Equal("123456789"))
		Expect(up.Email).Should(Equal("johndoe@example.com"))
		Expect(up.CPF).Should(Equal("12345678900"))
	})
	It("UpdateUser should fail", func() {
		user := User{
			ID:    1,
			Name:  "John",
			Age:   30,
			Phone: "123456789",
			Email: "johndoe@example.com",
			CPF:   "12345678900",
		}
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE Users SET Name = $1, Age = $2, Phone = $3, Email = $4, CPF = $5 WHERE ID =$6`)).WithArgs(user.Name, user.Age, user.Phone, user.Email, user.CPF, user.ID).WillReturnError(errors.New("error while UpdateUser"))
		up, err := rep.UpdateUser(context.Background(), user)
		Expect(err).Should(HaveOccurred())
		Expect(up).Should(BeNil())
	})
	It("should DeleteUser successfully", func() {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM Users WHERE ID = ?`)).WithArgs(1).WillReturnResult(pgxmock.NewResult("", 1))
		err := rep.DeleteUser(context.Background(), 1)
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("DeleteUser should fail", func() {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM Users WHERE ID = ?`)).WithArgs(1).WillReturnError(errors.New("error while DeleteUser"))
		err := rep.DeleteUser(context.Background(), 1)
		Expect(err).Should(HaveOccurred())
	})
	It("should DeletePixKey successfully", func() {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM PixKey WHERE ID = ?`)).WithArgs("1").WillReturnResult(pgxmock.NewResult("", 1))
		err := rep.DeletePixKey(context.Background(), "1")
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("DeletePixKey should fail", func() {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM PixKey WHERE ID = ?`)).WithArgs("1").WillReturnError(errors.New("error while DeletePixKey"))
		err := rep.DeletePixKey(context.Background(), "1")
		Expect(err).Should(HaveOccurred())
	})
})
