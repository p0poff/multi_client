package main

import (
	"encoding/json"
	"fmt"
)

var rs Requests
var res Responses

func start() {
	rs = Requests{}
	res = Responses{}
}

func addRequest(url string, method string, body string, headers_json string) bool {
	headers := map[string]string{}
	err := json.Unmarshal([]byte(headers_json), &headers)

	if err != nil {
		return false
	}

	r := Request{
		url:     url,
		method:  method,
		body:    body,
		headers: headers,
	}
	rs.add(r)

	return true
}

func main() {
	fmt.Println("Hello, World!")

	start()
	addRequest("http://localhost:8080", "GET", "1", `{"Content-Type": "application/json"}`)
	addRequest("http://localhost:8080", "GET", "2", `{"Content-Type": "application/json"}`)

	fmt.Println(rs)

	start()
	addRequest("http://localhost:8080", "GET", "pop", `{"Content-Type": "application/json"}`)
	fmt.Println(rs)
}
