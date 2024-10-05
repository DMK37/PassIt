package api

import "github.com/DMK37/PassIt/feedservice/db"

type LikeRequest struct {
	OwnerId string `json:"owner_id"`
	PostId string `json:"post_id"`
}

type ResponsePost struct {
	Id        string   `json:"id" dynamodbav:"postId"`
	UserId    string   `json:"user_id" dynamodbav:"userId"`
	User      PostUser `json:"user" dynamodbav:"user"`
	Text      string   `json:"text" dynamodbav:"text"`
	Images    []string `json:"images" dynamodbav:"images"`
	Timestamp int64    `json:"timestamp" dynamodbav:"timestamp"`
	Likes     []string `json:"likes" dynamodbav:"likes"`
	Comments  []string `json:"comments" dynamodbav:"comments"`
}

type PostUser struct {
	Username  string `json:"username" dynamodbav:"username"`
	FirstName string `json:"first_name" dynamodbav:"first_name"`
	LastName  string `json:"last_name" dynamodbav:"last_name"`
	Avatar    string `json:"avatar" dynamodbav:"avatar"`
}

func mapPostToResponsePost(post *db.Post, user *db.User) *ResponsePost {
	postUser := PostUser{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	}
	return &ResponsePost{
		Id:        post.Id,
		UserId:    post.UserId,
		User:      postUser,
		Text:      post.Text,
		Images:    post.Images,
		Timestamp: post.Timestamp,
		Likes:     post.Likes,
		Comments:  post.Comments,
	}
}