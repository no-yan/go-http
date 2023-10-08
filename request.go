package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime"
	"strconv"
	"strings"
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

func (s *Server) parseRequest(conn net.Conn) (*Request, error) {
	r := bufio.NewReader(conn)

	requestLine, err := r.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("Connection closed: %w", err)
		}
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
		req.Body = string(body)
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

func PrintStack() {
	var pc [100]uintptr
	n := runtime.Callers(0, pc[:])
	frames := runtime.CallersFrames(pc[:n])
	var (
		fr runtime.Frame
		ok bool
	)
	if _, ok = frames.Next(); !ok {
		return
	}
	for {
		fr, ok = frames.Next()
		if !ok {
			return
		}
		fmt.Printf("%s:%d %s\n", fr.File, fr.Line, fr.Function)
	}
}
