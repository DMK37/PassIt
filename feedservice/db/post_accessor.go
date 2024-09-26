package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PostAccessor interface {
	CreatePost(post *Post) error
	GetPost(userId, postId string) (*Post, error)
	GetPosts(userId string, limit int32) ([]*Post, error)
}

type DynamoDBPostAccessor struct {
	db *dynamodb.Client
}

func NewDynamoDBPostAccessor() (*DynamoDBPostAccessor, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load config for production DynamoDB: %v", err)
	}
	db := dynamodb.NewFromConfig(cfg)

	return &DynamoDBPostAccessor{db: db}, nil
}

func (d *DynamoDBPostAccessor) CreatePost(post *Post) error {

	av, err := attributevalue.MarshalMap(post)
	if err != nil {
		return fmt.Errorf("failed to marshal post: %v", err)
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("PassItPosts"),
		Item:      av,
	}

	_, err = d.db.PutItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to put item: %v", err)
	}

	return nil
}

func (d *DynamoDBPostAccessor) GetPost(userId, postId string) (*Post, error) {

	input := &dynamodb.GetItemInput{
		TableName: aws.String("PassItPosts"),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userId},
			"postId": &types.AttributeValueMemberS{Value: postId},
		},
	}

	result, err := d.db.GetItem(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	post := &Post{}
	err = attributevalue.UnmarshalMap(result.Item, post)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal post: %v", err)
	}

	return post, nil
}

func (d *DynamoDBPostAccessor) GetPosts(userId string, limit int32) ([]*Post, error) {

	input := &dynamodb.QueryInput{
		TableName:              aws.String("PassItPosts"),
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userId},
		},
		Limit: aws.Int32(limit),
	}

	result, err := d.db.Query(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}

	posts := []*Post{}

	if err = attributevalue.UnmarshalListOfMaps(result.Items, &posts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal posts: %v", err)
	}

	return posts, nil
}
