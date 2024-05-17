package bankaccount

import (
	"context"
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
	return args.Error(1)
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
			//idTransaction := ksuid.New().String()
			//payerMovement := Movement{
			//	Amount:        "-20.00",
			//	UserID:        "1",
			//	Date:          "2024-05-17T10:22:26-03:00",
			//	TransactionID: idTransaction,
			//	OperationType: Debit,
			//}
			//receiverMovement := Movement{
			//	Amount:        "20.00",
			//	UserID:        "2",
			//	Date:          "2024-05-17T10:22:26-03:00",
			//	TransactionID: idTransaction,
			//	OperationType: Credit,
			//}
			repMock.On("Save", context.Background(), mock.MatchedBy(func(mov Movement) bool {
				Expect(mov.Amount).ShouldNot(BeEmpty())
				Expect(mov.OperationType).Should(Equal(Debit))
				Expect(mov.UserID).Should(Equal("1"))
				return true
			})).Return(nil)
			repMock.On("Save", context.Background(), mock.MatchedBy(func(mov Movement) bool {
				Expect(mov.Amount).ShouldNot(BeEmpty())
				Expect(mov.Amount).Should(Equal("20.00"))
				Expect(mov.OperationType).Should(Equal(Credit))
				Expect(mov.UserID).Should(Equal("2"))
				return true
			})).Return(nil)
			err := serv.TransferPIX(context.Background(), request)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
