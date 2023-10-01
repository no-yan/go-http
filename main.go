package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

type Request struct {
	Method        string
	Target        string
	ProtoVersion  string
	ContentLength int
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
	ln, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := ln.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is running at localhost:8888")

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			defer conn.Close()
			var req Request

			// Read all request
			r := bufio.NewReader(conn)
			requestLine, err := r.ReadString('\n')
			if err != nil {
				panic("failed to read request line")
			}
			requestLine = strings.TrimSpace(requestLine) // Remove trailing newline
			method, target, protoVersion, err := parseRequestLine(requestLine)
			if err != nil {
				panic(err)
			}
			req.Method = method
			req.Target = target
			req.ProtoVersion = protoVersion
			fmt.Printf("method=%s, target=%s, protoVersion=%s\n", method, target, protoVersion)
			// if protoVersion != "HTTP/1.1" {
			// 	fmt.Errorf("unsupported protocol version: %s", protoVersion)
			// }

			// read field lines
			for {
				line, err := r.ReadString('\n')
				if err != nil {
					panic("failed to read field lines")
				}
				line = strings.TrimSpace(line)

				if line == "" {
					fmt.Println("end of field lines")
					break // end of field lines
				}
				if strings.HasPrefix(line, "Content-Length:") {
					fmt.Sscanf(line, "Content-Length: %d", &req.ContentLength)
				}
			}

			// read body
			// TODO: read body only if method can have body
			if req.ContentLength > 0 {
				body := make([]byte, req.ContentLength)
				if _, err := r.Read(body); err != nil && err != io.EOF {
					panic(err)
				}
			}
			res.Write(conn)
		}()
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
