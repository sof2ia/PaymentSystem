package internal

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBClient implements dynamo necessary methods
type DynamoDBClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

// PutItem saves item
func PutItem(ctx context.Context, api DynamoDBClient, params *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return api.PutItem(ctx, params)
}

// Query lists item
func Query(ctx context.Context, api DynamoDBClient, params *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return api.Query(ctx, params)
}
