package internal

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

type mockDynamoDB struct {
	mock.Mock
	DynamoDBClient
}

func (m *mockDynamoDB) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *mockDynamoDB) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

var _ = Describe("Repository Test", Ordered, func() {
	Context("Save", func() {
		var rep Repository
		var clientDB *mockDynamoDB
		BeforeEach(func() {
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
			clientDB.On("PutItem", context.Background(), &dynamodb.PutItemInput{Item: movMap, TableName: aws.String("Movement")}, mock.Anything).Return(nil, errors.New("fields should not be empty"))
			err = rep.Save(context.Background(), movement)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("fields should not be empty"))
		})
	})
	Context("ListMovementsByUser", func() {
		var rep Repository
		var clientDB *mockDynamoDB
		BeforeEach(func() {
			clientDB = new(mockDynamoDB)
			rep = NewRepository(clientDB)
		})
		It("should List movements successfully", func() {
			idUser := "1"
			keyEx := expression.Key("UserID").Equal(expression.Value(idUser))
			expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
			if err != nil {
				panic(err)
				return
			}
			clientDB.On("Query", context.Background(), &dynamodb.QueryInput{TableName: aws.String("Movement"), ExpressionAttributeNames: expr.Names(), ExpressionAttributeValues: expr.Values(), KeyConditionExpression: expr.KeyCondition()}, mock.Anything).Return(&dynamodb.QueryOutput{
				ConsumedCapacity: nil,
				Items: []map[string]ddbTypes.AttributeValue{
					{
						"ID":     &ddbTypes.AttributeValueMemberS{Value: "1"},
						"Amount": &ddbTypes.AttributeValueMemberS{Value: "10.00"},
					},
					{
						"ID":     &ddbTypes.AttributeValueMemberS{Value: "2"},
						"Amount": &ddbTypes.AttributeValueMemberS{Value: "20.00"},
					},
					{
						"ID":     &ddbTypes.AttributeValueMemberS{Value: "3"},
						"Amount": &ddbTypes.AttributeValueMemberS{Value: "30.00"},
					},
				},
			}, nil)
			output, err := rep.ListMovementsByUser(context.Background(), "1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(output[0].ID).Should(Equal("1"))
			Expect(output[1].Amount).Should(Equal("20.00"))
			Expect(output[2].Amount).Should(Equal("30.00"))
		})
		It("should list movements unsuccessfully", func() {
			idUser := "2"

			keyEx := expression.Key("UserID").Equal(expression.Value(idUser))
			expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
			if err != nil {
				panic(err)
				return
			}
			clientDB.On("Query", context.Background(), &dynamodb.QueryInput{TableName: aws.String("Movement"), ExpressionAttributeNames: expr.Names(), ExpressionAttributeValues: expr.Values(), KeyConditionExpression: expr.KeyCondition()}, mock.Anything).Return(nil, errors.New("no movement found for user"))
			output, err := rep.ListMovementsByUser(context.Background(), "2")
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("no movement found for user"))
			Expect(len(output)).Should(Equal(0))
		})
	})
})
