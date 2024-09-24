package db

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	transport "github.com/aws/smithy-go/endpoints"
)

type UserAccessor interface {
	CreateUser(user *User) error
	GetUserById(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
}

type DynamoDBUserAccessor struct {
	db *dynamodb.Client
}

func NewDynamoDBUserAccessor() (*DynamoDBUserAccessor, error) {

	var cfg aws.Config
	var err error

	dynamoEndpoint := os.Getenv("DYNAMODB_ENDPOINT")

	if dynamoEndpoint != "" {

		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("us-east-1"),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load config for local DynamoDB: %v", err)
		}
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to load config for production DynamoDB: %v", err)
		}
	}

	fmt.Println("DynamoDB endpoint: ", dynamoEndpoint)
	db := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if dynamoEndpoint != "" {
			o.EndpointResolverV2 = &resolverV2{endpoint: dynamoEndpoint}
		}
	})

	return &DynamoDBUserAccessor{db: db}, nil
}

func (d *DynamoDBUserAccessor) CreateUser(user *User) error {
    av, err := attributevalue.MarshalMap(user)
    if err != nil {
        return fmt.Errorf("failed to marshal user: %v", err)
    }

    input := &dynamodb.PutItemInput{
        TableName: aws.String("PassItUsers"),
        Item:      av,
    }

	_, err = d.db.PutItem(context.TODO(), input)
    if err != nil {
        return fmt.Errorf("failed to put item: %v", err)
    }

    return nil
}

func (d *DynamoDBUserAccessor) GetUserById(id string) (*User, error) {
    input := &dynamodb.GetItemInput{
        TableName: aws.String("PassItUsers"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: id},
        },
    }

    result, err := d.db.GetItem(context.TODO(), input)
    if err != nil {
        return nil, fmt.Errorf("failed to get item: %v", err)
    }

    if result.Item == nil {
        return nil, nil
    }

    user := &User{}
    err = attributevalue.UnmarshalMap(result.Item, user)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal user: %v", err)
    }

    return user, nil
}

func (d *DynamoDBUserAccessor) GetUserByUsername(username string) (*User, error) {
    input := &dynamodb.QueryInput{
        TableName:              aws.String("PassItUsers"),
        IndexName:              aws.String("username-index"),
        KeyConditionExpression: aws.String("username = :username"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":username": &types.AttributeValueMemberS{Value: username},
        },
    }

    result, err := d.db.Query(context.TODO(), input)
    if err != nil {
        return nil, fmt.Errorf("failed to query user by username: %v", err)
    }

    if len(result.Items) == 0 {
        return nil, nil // No user found
    }

    user := &User{}
    err = attributevalue.UnmarshalMap(result.Items[0], user)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal user: %v", err)
    }

	

    return user, nil
}
type resolverV2 struct {
	endpoint string
}

func (r *resolverV2) ResolveEndpoint(ctx context.Context, params dynamodb.EndpointParameters) (transport.Endpoint, error) {

	u, err := url.Parse(r.endpoint)
	if err != nil {
		return transport.Endpoint{}, err
	}

	return transport.Endpoint{
		URI: *u,
	}, nil
}