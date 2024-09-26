package main

import "github.com/DMK37/PassIt/feedservice/api"

func main() {
	s := api.NewServer(":8081")
	s.Start()
}