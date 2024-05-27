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
	})
	Context("DepositAmount", func() {
		var serv Service
		var repMock *mockRepository
		BeforeEach(func() {
			repMock = new(mockRepository)
			serv = NewService(repMock)
		})
		It("should pass DepositAmount successfully", func() {
			deposit := DepositAmountRequest{
				Amount: 2000.00,
				UserID: "1",
			}
			repMock.On("Save", mock.Anything, mock.AnythingOfType("Movement")).Return(nil)

			err := serv.DepositAmount(context.Background(), deposit)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should pass DepositAmount unsuccessfully", func() {
			deposit := DepositAmountRequest{
				Amount: 2000.00,
				UserID: "1",
			}
			repMock.On("Save", mock.Anything, mock.AnythingOfType("Movement")).Return(errors.New("error while DepositAmount"))
			err := serv.DepositAmount(context.Background(), deposit)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("error while DepositAmount"))
		})
	})
})
