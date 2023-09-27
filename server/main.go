package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

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
		go func() {
			defer conn.Close()
			// Read all request
			b := make([]byte, 4096)
			n, err := conn.Read(b)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					panic(err)
				}
			}
			fmt.Println(string(b[:n]))
		}()
	}
}
