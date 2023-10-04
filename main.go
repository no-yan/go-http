package main

import (
	"io"
	"log"
)

const (
	ProtocolVersion = "HTTP/1.1"
	ServerAddress   = "localhost:8888"
)

type Request struct {
	Method        string
	Target        string
	ProtoVersion  string
	ContentLength int
	Fields        map[string][]string
}

type Response struct {
}

func (*Response) Write(w io.Writer) {
	w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	w.Write([]byte("Content-Length: 2\r\n"))
	w.Write([]byte("\r\n"))
	w.Write([]byte("OK"))
}

func main() {
	server := NewServer(ServerAddress)
	if err := server.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
