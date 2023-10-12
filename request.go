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

	if err := req.parseRequestLine(r); err != nil {
		return fmt.Errorf("failed to parse request line: %v", err)
	}

	if err := parseHeader(r, req); err != nil {
		return fmt.Errorf("failed to parse header: %v", err)
	}

	if err := req.parseBody(r); err != nil {
		return fmt.Errorf("failed to parse body: %v", err)
	}
	return nil
}

// Parse request line ("GET /PATH HTTP/1.1") to three parts
func (req *Request) parseRequestLine(r *bufio.Reader) error {
	line, err := r.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("Connection closed: %w", err)
		}
		return err
	}
	if req.Method, req.Target, req.ProtoVersion, err = parseRequestLine(line); err != nil {
		return err
	}

	if req.ProtoVersion != "HTTP/1.1" {
		// Deny request
		return fmt.Errorf("unsupported protocol version: %s", req.ProtoVersion)
	}

	return nil
}

func parseRequestLine(line string) (method, target, protoVersion string, err error) {
	line = strings.TrimSpace(line) // Remove trailing newline
	ls := strings.Split(line, " ")
	if len(ls) != 3 {
		err = fmt.Errorf("invalid request line: %s", line)
	}

	method = ls[0]
	target = ls[1]
	protoVersion = ls[2]
	return
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

func (req *Request) parseBody(r *bufio.Reader) error {
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
