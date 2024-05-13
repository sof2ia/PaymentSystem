package bankaccount

import (
	"PaymentSystem/bankaccount/internal"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

type mockDynamoDB struct {
	mock.Mock
	internal.DynamoDBClient
}

func (m *mockDynamoDB) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

var _ = Describe("Repository Test", Ordered, func() {
	Context("Save", func() {
		var rep Repository
		var clientDB *mockDynamoDB
		BeforeAll(func() {
			clientDB = new(mockDynamoDB)
			rep = NewRepository(clientDB)
		})
		It("should Save successfully", func() {
			movement := Movement{
				ID:            "sdf56yt",
				Amount:        "20",
				UserID:        "1",
				Date:          "2024-05-05",
				TransactionID: "d6t",
				OperationType: "credit",
			}
			movMap, err := attributevalue.MarshalMap(movement)
			if err != nil {
				panic(err)
				return
			}
			clientDB.On("PutItem", context.Background(), &dynamodb.PutItemInput{Item: movMap, TableName: aws.String("Movement")}, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)
			err = rep.Save(context.Background(), movement)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should Save unsuccessfully", func() {
			movement := Movement{
				ID:            "sdf56yt",
				Amount:        "20",
				UserID:        "1",
				Date:          "2024-05-05",
				TransactionID: "d6t",
				OperationType: "credit",
			}
			movMap, err := attributevalue.MarshalMap(movement)
			clientDB.On("PutItem", context.Background(), &dynamodb.PutItemInput{Item: movMap, TableName: aws.String("Movement")}, mock.Anything).Return(&dynamodb.PutItemOutput{}, errors.New("fields should not be empty"))
			err = rep.Save(context.Background(), movement)
			Expect(err).Should(HaveOccurred())
		})
	})
})
