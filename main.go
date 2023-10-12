package main

import (
	"log"
)

const (
	ProtocolVersion = "HTTP/1.1"
	ServerAddress   = "localhost:8888"
)

func main() {
	server := NewServer(ServerAddress)
	if err := server.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
