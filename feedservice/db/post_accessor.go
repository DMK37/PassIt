package db

import (
	"context"
	"fmt"
	"sort"

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
	GetFollowingPosts(userId string) ([]*Post, error)
	GetPostUser(userId string) (*PostUser, error)
	LikePost(userId, postId, ownerId string) error
	UnlikePost(userId, postId, ownerId string) error
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

func (d *DynamoDBPostAccessor) GetPostUser(userId string) (*PostUser, error) {
	inputU := &dynamodb.GetItemInput{
		TableName: aws.String("PassItUsers"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: userId},
		},
	}

	result, err := d.db.GetItem(context.TODO(), inputU)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %v", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("user not found")
	}

	user := &PostUser{}
	err = attributevalue.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	return user, nil
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

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Timestamp > posts[j].Timestamp
	})

	return posts, nil
}

func (d *DynamoDBPostAccessor) GetFollowingPosts(userId string) ([]*Post, error) {

	input := &dynamodb.GetItemInput{
		TableName: aws.String("PassItUsers"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: userId},
		},
	}

	result, err := d.db.GetItem(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %v", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("user not found")
	}

	user := &User{}
	err = attributevalue.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	following := user.Following

	res := []*Post{}

	for _, followingId := range following {
		input := &dynamodb.QueryInput{
			TableName:              aws.String("PassItPosts"),
			KeyConditionExpression: aws.String("userId = :userId"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":userId": &types.AttributeValueMemberS{Value: followingId},
			},
		}

		result, err := d.db.Query(context.TODO(), input)
		if err != nil {
			return nil, fmt.Errorf("failed to query: %v", err)
		}

		posts := []*Post{}
		if err = attributevalue.UnmarshalListOfMaps(result.Items, &posts); err != nil {
			return nil, fmt.Errorf("failed to unmarshal posts: %v", err)
		}

		res = append(res, posts...)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Timestamp > res[j].Timestamp
	})

	return res, nil
}

func (d *DynamoDBPostAccessor) LikePost(userId, postId, OwnerId string) error {

	post, err := d.GetPost(OwnerId, postId)
	if err != nil {
		return fmt.Errorf("failed to get post: %v", err)
	}

	if post == nil {
		return fmt.Errorf("post not found")
	}

	if !contains(post.Likes, userId) {
		post.Likes = append(post.Likes, userId)
	}

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

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func (d *DynamoDBPostAccessor) UnlikePost(userId, postId, OwnerId string) error {

	post, err := d.GetPost(OwnerId, postId)
	if err != nil {
		return fmt.Errorf("failed to get post: %v", err)
	}

	if post == nil {
		return fmt.Errorf("post not found")
	}

	if contains(post.Likes, userId) {
		post.Likes = remove(post.Likes, userId)
	}

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

func remove(arr []string, str string) []string {

	index := -1
	for i, a := range arr {
		if a == str {
			index = i
			break
		}
	}

	if index == -1 {
		return arr
	}

	arr[index] = arr[len(arr)-1]
	return arr[:len(arr)-1]
}
