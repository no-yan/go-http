package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
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

	if protoVersion != "HTTP/1.1" {
		// Deny request
		return nil, fmt.Errorf("unsupported protocol version: %s", protoVersion)
	}

	req := &Request{
		Method:       method,
		Target:       target,
		ProtoVersion: protoVersion,
		Fields:       make(map[string][]string),
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
			return nil, fmt.Errorf("field line should have separator ':' || %s", line)
		}

		if name == "Content-Length" {
			req.ContentLength, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid content length: %v", err)
			}
		}

		if _, ok := req.Fields[name]; !ok {
			req.Fields[name] = []string{}
		}
		req.Fields[name] = append(req.Fields[name], value)
	}

	// read body
	// TODO: read body only if method can have body
	if req.ContentLength > 0 {
		body := make([]byte, req.ContentLength)
		if _, err := r.Read(body); err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read body: %v", err)
		}
	}

	return req, nil
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
