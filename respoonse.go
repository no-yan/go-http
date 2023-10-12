package main

import (
	"fmt"
	"io"
)

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
