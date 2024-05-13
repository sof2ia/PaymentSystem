package bankaccount

import (
	"PaymentSystem/bankaccount/internal"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
)

type Repository struct {
	Client internal.DynamoDBClient
}

func (r *Repository) Save(ctx context.Context, mov Movement) error {
	movMap, err := attributevalue.MarshalMap(mov)
	if err != nil {
		return err
	}
	putItem := &dynamodb.PutItemInput{
		Item:      movMap,
		TableName: aws.String("Movement"),
	}
	log.Print(putItem)
	_, err = internal.PutItem(ctx, r.Client, putItem)
	log.Print(err)
	return err
}

func NewRepository(client internal.DynamoDBClient) Repository {
	return Repository{Client: client}
}
