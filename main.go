package main

import (
	"fmt"
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
	Response      *Response
	Body          string
}

type Response struct {
	StatusCode    int
	StatusMessage string
	Fields        map[string][]string
	Body          string
	ProtoVersion  string
	ContentLength int
	Request       *Request
}

func (r *Response) Write(w io.Writer) {
	fmt.Fprintf(w, "%s %d %s\r\n", r.ProtoVersion, r.StatusCode, r.StatusMessage)
	fmt.Fprintf(w, "Content-Length: %d\r\n", len(r.Body))
	fmt.Fprintf(w, "\r\n")
	fmt.Fprintf(w, "%s", r.Body)
}

func main() {
	server := NewServer(ServerAddress)
	if err := server.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
