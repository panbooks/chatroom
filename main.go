package main

import "chat/server"

func main() {
	server := server.NewServer()
	server.Start()
}
