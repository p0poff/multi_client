package main

import (
	"fmt"
)

var rs Requests
var res Responses

func start() {
	rs = Requests{}
	res = Responses{}
}

func main() {
	fmt.Println("Hello, World!")
	r := Request{
		url:    "http://localhost:8080",
		method: "GET",
		body:   "",
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	start()
	rs.add(r)
	rs.add(r)
	rs.add(r)

	fmt.Println(rs)

	start()
	rs.add(r)
	fmt.Println(rs)
}
