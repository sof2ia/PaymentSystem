package internal

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Repository interface {
	Save(ctx context.Context, mov Movement) error
	ListMovementsByUser(ctx context.Context, idUser string) ([]Movement, error)
}

type repository struct {
	Client DynamoDBClient
}

func (r *repository) Save(ctx context.Context, mov Movement) error {
	movMap, err := attributevalue.MarshalMap(mov)
	if err != nil {
		return err
	}
	putItem := &dynamodb.PutItemInput{
		Item:      movMap,
		TableName: aws.String("Movement"),
	}
	_, err = PutItem(ctx, r.Client, putItem)
	return err
}

func (r *repository) ListMovementsByUser(ctx context.Context, idUser string) ([]Movement, error) {
	keyEx := expression.Key("UserID").Equal(expression.Value(idUser))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}
	input := &dynamodb.QueryInput{
		TableName:                 aws.String("Movement"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}
	output, err := Query(ctx, r.Client, input)
	if err != nil {
		return nil, err
	}
	var movements []Movement
	err = attributevalue.UnmarshalListOfMaps(output.Items, &movements)
	if err != nil {
		return nil, err
	}
	return movements, err
}

func NewRepository(client DynamoDBClient) Repository {
	return &repository{Client: client}
}
