package db

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id        string   `json:"id" dynamodbav:"postId"`
	UserId    string   `json:"user_id" dynamodbav:"userId"`
	User      PostUser `json:"user" dynamodbav:"user"`
	Text      string   `json:"text" dynamodbav:"text"`
	Images    []string `json:"images" dynamodbav:"images"`
	Timestamp int64    `json:"timestamp" dynamodbav:"timestamp"`
	Likes     []string `json:"likes" dynamodbav:"likes"`
}

type PostUser struct {
	Username  string `json:"username" dynamodbav:"username"`
	FirstName string `json:"first_name" dynamodbav:"first_name"`
	LastName  string `json:"last_name" dynamodbav:"last_name"`
	Avatar    string `json:"avatar" dynamodbav:"avatar"`
}

func NewPost(userId string, text string, images []string, user PostUser) *Post {
	return &Post{
		Id:        uuid.New().String(),
		UserId:    userId,
		User:      user,
		Text:      text,
		Images:    images,
		Timestamp: time.Now().Unix(),
		Likes:     []string{},
	}
}
