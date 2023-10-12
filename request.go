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

func NewRequest() *Request {
	return &Request{
		Fields: map[string][]string{},
	}
}

func (req *Request) parseRequest(conn net.Conn) error {
	r := bufio.NewReader(conn)

	requestLine, err := r.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("Connection closed: %w", err)
		}
		return err
	}
	method, target, protoVersion, err := parseRequestLine(requestLine)
	if err != nil {
		return err
	}

	if protoVersion != "HTTP/1.1" {
		// Deny request
		return fmt.Errorf("unsupported protocol version: %s", protoVersion)
	}

	req.Method = method
	req.Target = target
	req.ProtoVersion = protoVersion

	if err := parseHeader(r, req); err != nil {
		return fmt.Errorf("failed to parse header: %v", err)
	}

	if err := parseBody(r, req); err != nil {
		return fmt.Errorf("failed to parse body: %v", err)
	}
	return nil
}

// Parse request line ("GET /PATH HTTP/1.1") to three parts
func parseRequestLine(line string) (string, string, string, error) {
	line = strings.TrimSpace(line) // Remove trailing newline
	ls := strings.Split(line, " ")

	if len(ls) != 3 {
		return "", "", "", fmt.Errorf("invalid request line: %s", line)
	}
	method := ls[0]
	target := ls[1]
	protoVersion := ls[2]

	return method, target, protoVersion, nil
}

func parseHeader(r *bufio.Reader, req *Request) error {
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read header line: %v", err)
		}
		line = strings.TrimSpace(line)

		if line == "" {
			break // end of field lines
		}

		name, value, found := strings.Cut(line, ":")
		name, value = strings.TrimSpace(name), strings.TrimSpace(value)
		if !found {
			return fmt.Errorf("field line should have separator ':' || %s", line)
		}

		if name == "Content-Length" {
			req.ContentLength, err = strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid content length: %v", err)
			}
		}

		if _, ok := req.Fields[name]; !ok {
			req.Fields[name] = []string{}
		}
		req.Fields[name] = append(req.Fields[name], value)
	}
	return nil
}

func parseBody(r *bufio.Reader, req *Request) error {
	// TODO: read body only if method can have body
	if req.ContentLength > 0 {
		body := make([]byte, req.ContentLength)
		if _, err := r.Read(body); err != nil && err != io.EOF {
			return fmt.Errorf("failed to read body: %v", err)
		}
		req.Body = string(body)
	}
	return nil
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
