package main

import (
	"fmt"
)

type request struct {
	url     string
	method  string
	body    string
	headers map[string]string
}

type requests []request

func (r *requests) add(r2 request) {
	*r = append(*r, r2)
}

func main() {
	fmt.Println("Hello, World!")
	r := request{
		url:    "http://localhost:8080",
		method: "GET",
		body:   "",
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	rs := requests{}
	rs.add(r)
	rs.add(r)
	rs.add(r)

	fmt.Println(rs)
}
