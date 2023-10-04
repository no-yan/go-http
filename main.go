package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	ProtocolVersion = "HTTP/1.1"
	ServerAddress   = "localhost:8888"
)

type Server struct {
	Address string
}

func NewServer(address string) *Server {
	return &Server{Address: address}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	defer func() {
		err := ln.Close()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println("Server is running at localhost:8888")

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := s.parseRequest(conn)
	if err != nil {
		s.sendErrorResponse(conn, err)
		return
	}

	s.processRequest(conn, req)
}

func (s *Server) parseRequest(conn net.Conn) (*Request, error) {
	r := bufio.NewReader(conn)

	requestLine, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	requestLine = strings.TrimSpace(requestLine) // Remove trailing newline
	method, target, protoVersion, err := parseRequestLine(requestLine)
	if err != nil {
		return nil, err
	}

	req := &Request{
		Method:       method,
		Target:       target,
		ProtoVersion: protoVersion,
		Fields:       make(map[string]string),
	}

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read header line: %v", err)
		}
		line = strings.TrimSpace(line)

		if line == "" {
			break // end of field lines
		}

		name, value, found := strings.Cut(line, ":")
		name, value = strings.TrimSpace(name), strings.TrimSpace(value)
		if !found {
			// フィールドラインに ":" が含まれていないというエラーを返す
			return nil, fmt.Errorf("field line should have separator ':' || %s", line)

		}

		if name == "Content-Length" {
			req.ContentLength, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid content length: %v", err)
			}
		}

		if _, ok := req.Fields[name]; !ok {
			req.Fields[name] = strings.TrimSpace(value)
		}

	}

	return req, nil
}

func (s *Server) sendErrorResponse(conn net.Conn, err error) {
	// ...
}

func (s *Server) processRequest(conn net.Conn, req *Request) {
	// ...
}

type Request struct {
	Method        string
	Target        string
	ProtoVersion  string
	ContentLength int
	Fields        map[string]string
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

// Parse request line ("GET /PATH HTTP/1.1") to three parts
func parseRequestLine(line string) (string, string, string, error) {
	ls := strings.Split(line, " ")

	if len(ls) != 3 {
		return "", "", "", fmt.Errorf("invalid request line: %s", line)
	}
	method := ls[0]
	target := ls[1]
	protoVersion := ls[2]

	return method, target, protoVersion, nil
}
