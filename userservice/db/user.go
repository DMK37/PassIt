package db

import (
	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	Id        string   `json:"id" dynamodbav:"id"`
	Username  string   `json:"username" dynamodbav:"username"`
	Email     string   `json:"email" dynamodbav:"email"`
	Password  string   `json:"password" dynamodbav:"password"`
	FirstName string   `json:"first_name" dynamodbav:"first_name"`
	LastName  string   `json:"last_name" dynamodbav:"last_name"`
	Followers []string `json:"followers" dynamodbav:"followers"`
	Following []string `json:"following" dynamodbav:"following"`
	Avatar    string   `json:"avatar" dynamodbav:"avatar"`
}

func NewUser(username, email, password, firstName, lastName string) *User {
	return &User{
		Id:        uuid.New().String(),
		Username:  username,
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Followers: []string{},
		Following: []string{},
		Avatar:    "",
	}
}

func (u *User) String() string {
	jsonData, err := json.Marshal(u)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
