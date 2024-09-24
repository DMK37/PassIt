package main

import (
	"userservice/api"
)

func main() {
	s := api.NewServer(":8080")
	s.Start()
}