package internal

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBClient implements dynamo necessary methods
type DynamoDBClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

// PutItem saves item
func PutItem(ctx context.Context, api DynamoDBClient, params *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return api.PutItem(ctx, params)
}
