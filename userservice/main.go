package main

import "github.com/DMK37/PassIt/userservice/api"

func main() {
	s := api.NewServer(":8080")
	s.Start()
}