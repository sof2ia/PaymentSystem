package bankaccount

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Repository Test", func() {
	Context("Save", func() {
		var rep Repository
		var db dynamodb.Client
		BeforeAll(func() {
			rep = NewRepository(db)
		})
		It("should Save successfully", func() {
			rep.Save(context.Background(), Movement{
				ID:            "sdf56yt",
				Amount:        "20",
				UserID:        "1",
				Date:          "2024-05-05",
				TransactionID: "d6t",
				OperationType: "credit",
			})
		})
	})
})
