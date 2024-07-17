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
			repMock.On("ListMovementsByUser", mock.Anything, request.PayerID).Return([]Movement{
				{
					ID:            "1",
					Amount:        "100",
					Date:          "1/07/2024",
					TransactionID: "rhk5",
					OperationType: Credit,
				},
			}, nil)
			repMock.On("Save", context.Background(), mock.AnythingOfType("Movement")).Return(nil)
			repMock.On("Save", context.Background(), mock.AnythingOfType("Movement")).Return(nil)

			err := serv.TransferPIX(context.Background(), request)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should TransferPIX unsuccessfully - insufficient balance", func() {
			request := TransferRequest{
				PayerID:        "1",
				ReceiverPixKey: "2",
				Amount:         20.00,
			}
			repMock.On("ListMovementsByUser", mock.Anything, request.PayerID).Return([]Movement{
				{
					ID:            "1",
					Amount:        "100",
					Date:          "1/07/2024",
					TransactionID: "rhk5",
					OperationType: Credit,
				},
				{
					ID:            "2",
					Amount:        "-90",
					Date:          "1/07/2024",
					TransactionID: "udve4",
					OperationType: Debt,
				},
			}, errors.New("insufficient balance"))
			err := serv.TransferPIX(context.Background(), request)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("insufficient balance"))
		})
		It("should TransferPIX unsuccessfully while saving", func() {
			request := TransferRequest{
				PayerID:        "1",
				ReceiverPixKey: "2",
				Amount:         20.00,
			}
			repMock.On("ListMovementsByUser", mock.Anything, request.PayerID).Return([]Movement{
				{
					ID:            "1",
					Amount:        "100",
					Date:          "1/07/2024",
					TransactionID: "rhk5",
					OperationType: Credit,
				},
				{
					ID:            "2",
					Amount:        "-80",
					Date:          "1/07/2024",
					TransactionID: "udve4",
					OperationType: Debt,
				},
			}, nil)
			repMock.On("Save", context.Background(), mock.AnythingOfType("Movement")).Return(errors.New("error while saving"))

			err := serv.TransferPIX(context.Background(), request)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("error while saving"))
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
			repMock.On("ListMovementsByUser", mock.Anything, deposit.UserID).Return([]Movement{
				{
					ID:            "1",
					Amount:        "100",
					Date:          "1/07/2024",
					TransactionID: "rhk5",
					OperationType: Credit,
				},
				{
					ID:            "2",
					Amount:        "-80",
					Date:          "1/07/2024",
					TransactionID: "udve4",
					OperationType: Debt,
				},
			}, nil)
			repMock.On("Save", mock.Anything, mock.AnythingOfType("Movement")).Return(nil)

			err := serv.DepositAmount(context.Background(), deposit)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("DepositAmount give an error while listing MovementsByUser", func() {
			deposit := DepositAmountRequest{
				Amount: 2000.00,
				UserID: "1",
			}
			repMock.On("ListMovementsByUser", mock.Anything, deposit.UserID).Return(nil, errors.New("error in the List"))
			err := serv.DepositAmount(context.Background(), deposit)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("error in the List"))
		})
		It("should pass DepositAmount unsuccessfully while Saving", func() {
			deposit := DepositAmountRequest{
				Amount: 2000.00,
				UserID: "1",
			}
			repMock.On("ListMovementsByUser", mock.Anything, deposit.UserID).Return([]Movement{
				{
					ID:            "1",
					Amount:        "100",
					Date:          "1/07/2024",
					TransactionID: "rhk5",
					OperationType: Credit,
				},
				{
					ID:            "2",
					Amount:        "-80",
					Date:          "1/07/2024",
					TransactionID: "udve4",
					OperationType: Debt,
				},
			}, nil)
			repMock.On("Save", mock.Anything, mock.AnythingOfType("Movement")).Return(errors.New("error while saving"))

			err := serv.DepositAmount(context.Background(), deposit)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("error while saving"))
		})

	})
	Context("GetBalance Test", func() {
		var serv Service
		var repMock *mockRepository
		BeforeEach(func() {
			repMock = new(mockRepository)
			serv = NewService(repMock)
		})
		It("should GetBalance successfully", func() {
			repMock.On("ListMovementsByUser", mock.Anything, "1").Return([]Movement{
				{
					ID:            "1",
					Amount:        "100",
					Date:          "1/07/2024",
					TransactionID: "rhk5",
					OperationType: Credit,
				},
				{
					ID:            "2",
					Amount:        "-80",
					Date:          "1/07/2024",
					TransactionID: "udve4",
					OperationType: Debt,
				},
			}, nil)

			idRequest := GetBalanceRequest{ID: 1}
			res, err := serv.GetBalance(context.Background(), idRequest)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(res.Balance).Should(Equal(float64(20)))
		})
		It("should GetBalance unsuccessfully", func() {
			repMock.On("ListMovementsByUser", mock.Anything, "1").Return(([]Movement)(nil), errors.New("error while GetBalance"))
			idRequest := GetBalanceRequest{ID: 1}
			res, err := serv.GetBalance(context.Background(), idRequest)
			Expect(err).Should(HaveOccurred())
			Expect(res).Should(BeNil())
		})
	})
})
