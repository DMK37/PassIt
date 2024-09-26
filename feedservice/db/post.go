package db

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id     string   `json:"id" dynamodbav:"postId"`
	UserId string   `json:"user_id" dynamodbav:"userId"`
	Text   string   `json:"text" dynamodbav:"text"`
	Images []string `json:"images" dynamodbav:"images"`
	Timestamp int64 `json:"timestamp" dynamodbav:"timestamp"`
}


func NewPost(userId, text string, images []string) *Post {
	return &Post{
		Id:     uuid.New().String(),
		UserId: userId,
		Text:   text,
		Images: images,
		Timestamp: time.Now().Unix(),
	}
}