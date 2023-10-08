package main

import (
	"errors"
	"fmt"
	"io"
	"net"
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
		PrintStack()
		return err
	}
	defer ln.Close()

	fmt.Println("Server is running on", s.Address, "...")

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
		if errors.Is(err, io.EOF) {
			return // client closed connection
		}
		PrintStack()
		fmt.Println(err)
		s.sendErrorResponse(conn, err)
		return
	}

	s.processRequest(conn, req)
}

func (s *Server) sendErrorResponse(conn net.Conn, err error) {
	res := Response{
		StatusCode:    500,
		StatusMessage: "Internal Server Error",
		ProtoVersion:  "HTTP/1.1",
		Body:          err.Error(),
	}
	res.Write(conn)
}

func (s *Server) processRequest(conn net.Conn, req *Request) {
	r := Response{
		StatusCode:    200,
		StatusMessage: "OK",
		ProtoVersion:  "HTTP/1.1",
		Body:          "OK",
	}
	r.Write(conn)
}
