package db

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	transport "github.com/aws/smithy-go/endpoints"
)

type UserAccessor interface {
	CreateUser(user *User) error
	GetUserById(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	FollowUser(userId, followUserId string) error
	UnfollowUser(userId, followUserId string) error
	UpdateUser(userId, username, firstName, lastName, avatarPath string) error
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

func (d *DynamoDBUserAccessor) GetUserByEmail(email string) (*User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("PassItUsers"),
		IndexName:              aws.String("email-index"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := d.db.Query(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %v", err)
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

func (d *DynamoDBUserAccessor) FollowUser(userId, followUserId string) error {
	// Get user by ID
	user, err := d.GetUserById(userId)
	if err != nil {
		return fmt.Errorf("failed to get user by ID: %v", err)
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Get follow user by ID
	followUser, err := d.GetUserById(followUserId)
	if err != nil {
		return fmt.Errorf("failed to get follow user by ID: %v", err)
	}

	if followUser == nil {
		return fmt.Errorf("follow user not found")
	}

	// Add follow user ID to user's following list if not already following

	for _, followingId := range user.Following {
		if followingId == followUserId {
			return nil
		}
	}
	user.Following = append(user.Following, followUserId)

	// Add user ID to follow user's followers list
	followUser.Followers = append(followUser.Followers, userId)

	// Update user
	err = d.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	// Update follow user
	err = d.CreateUser(followUser)
	if err != nil {
		return fmt.Errorf("failed to update follow user: %v", err)
	}

	return nil
}

func (d *DynamoDBUserAccessor) UnfollowUser(userId, followUserId string) error {
	// Get user by ID
	user, err := d.GetUserById(userId)
	if err != nil {
		return fmt.Errorf("failed to get user by ID: %v", err)
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Get follow user by ID
	followUser, err := d.GetUserById(followUserId)
	if err != nil {
		return fmt.Errorf("failed to get follow user by ID: %v", err)
	}

	if followUser == nil {
		return fmt.Errorf("follow user not found")
	}

	// Remove follow user ID from user's following list
	for i, followingId := range user.Following {
		if followingId == followUserId {
			user.Following = append(user.Following[:i], user.Following[i+1:]...)
			break
		}
	}

	// Remove user ID from follow user's followers list
	for i, followerId := range followUser.Followers {
		if followerId == userId {
			followUser.Followers = append(followUser.Followers[:i], followUser.Followers[i+1:]...)
			break
		}
	}

	// Update user
	err = d.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	// Update follow user
	err = d.CreateUser(followUser)
	if err != nil {
		return fmt.Errorf("failed to update follow user: %v", err)
	}

	return nil
}

func (d *DynamoDBUserAccessor) UpdateUser(userId string, username string, firstName string, lastName string, avatarPath string) error {

	user, err := d.GetUserById(userId)
	if err != nil {
		return fmt.Errorf("failed to get user by ID: %v", err)
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.Username = username
	user.FirstName = firstName
	user.LastName = lastName
	if avatarPath != "" {
	user.Avatar = avatarPath
	}
	err = d.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
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
