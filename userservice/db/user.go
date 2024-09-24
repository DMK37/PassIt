package db

import "github.com/google/uuid"

type User struct {
	Id        string `json:"id" dynamodbav:"id"`
	Username  string `json:"username" dynamodbav:"username"`
	Password  string `json:"password" dynamodbav:"password"`
	FirstName string `json:"first_name" dynamodbav:"first_name"`
	LastName  string `json:"last_name" dynamodbav:"last_name"`
}

func NewUser(username, password, firstName, lastName string) *User {
	return &User{
		Id:        uuid.New().String(),
		Username:  username,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}
}
