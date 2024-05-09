package bankaccount

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Repository struct {
	db dynamodb.Client
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
	_, err = r.db.PutItem(ctx, putItem)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(db dynamodb.Client) Repository {
	return Repository{db}
}
