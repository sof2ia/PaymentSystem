package bankaccount

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

func (m *mockRepository) Save(ctx context.Context, mov Movement) error {
	args := m.Called(ctx, mov)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return args.Error(0)
}

func (m *mockRepository) ListMovementsByUser(ctx context.Context, idUser string) ([]Movement, error) {
	args := m.Called(ctx, idUser)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movement), args.Error(1)
}

var _ = Describe("Service Test", func() {
	Context("TransferPIX", func() {
		var serv Service
		var repMock *mockRepository
		BeforeEach(func() {
			repMock = new(mockRepository)
			serv = NewService(repMock)
		})
		It("should TransferPIX successfully", func() {
			request := TransferRequest{
				PayerID:        "1",
				ReceiverPixKey: "2",
				Amount:         20.00,
			}

			repMock.On("Save", context.Background(), mock.AnythingOfType("Movement")).Return(nil)
			repMock.On("Save", context.Background(), mock.AnythingOfType("Movement")).Return(nil)

			err := serv.TransferPIX(context.Background(), request)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should TransferPIX unsuccessfully", func() {
			request := TransferRequest{
				PayerID:        "1",
				ReceiverPixKey: "2",
				Amount:         20.00,
			}
			repMock.On("Save", context.Background(), mock.AnythingOfType("Movement")).Return(errors.New("error while Save payerMovement"))

			err := serv.TransferPIX(context.Background(), request)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("error while Save payerMovement"))
		})
		It("should fail validation while TransferPIX - empty fields", func() {
			request := TransferRequest{
				PayerID:        "",
				ReceiverPixKey: "",
				Amount:         20.00,
			}
			err := serv.TransferPIX(context.Background(), request)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("validation error: field: PayerID, value: \nvalidation error: field: ReceiverPixKey, value: \n"))
		})
		It("should fail validation while TransferPIX - amount == 0", func() {
			request := TransferRequest{
				PayerID:        "1",
				ReceiverPixKey: "2",
				Amount:         0.00,
			}
			err := serv.TransferPIX(context.Background(), request)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("validation error: field: Amount, value: 0.000000\n"))
		})
		It("should fail validation while TransferPIX - amount < 0", func() {
			request := TransferRequest{
				PayerID:        "1",
				ReceiverPixKey: "2",
				Amount:         -20.00,
			}
			err := serv.TransferPIX(context.Background(), request)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("validation error: field: Amount, value: -20.000000\n"))
		})
		It("should fail validation while TransferPIX - amount > 5000", func() {
			request := TransferRequest{
				PayerID:        "1",
				ReceiverPixKey: "2",
				Amount:         5001.00,
			}
			err := serv.TransferPIX(context.Background(), request)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("validation error: field: Amount, value: 5001.000000\n"))
		})
	})
})
