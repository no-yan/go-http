package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

type Request struct {
	Method       string
	Target       string
	ProtoVersion string
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
		var req Request
		go func() {
			defer conn.Close()
			// Read all request
			sc := bufio.NewScanner(conn)
			if sc.Scan() {
				if err := sc.Err(); err != nil {
					panic(err)
				}
				requestLine := sc.Text()
				method, target, protoVersion, err := parseRequestLine(requestLine)
				if err != nil {
					panic(err)
				}
				req.Method = method
				req.Target = target
				req.ProtoVersion = protoVersion
				if err != nil {
					panic(err)
				}
				fmt.Printf("method=%s, target=%s, protoVersion=%s\n", method, target, protoVersion)
			}

			res := Response{}
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
