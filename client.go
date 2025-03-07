package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

var rs Requests
var res Responses
var mu sync.Mutex

func start() {
	rs.clear()
	res.clear()
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

func send() {
	var wg sync.WaitGroup
	for _, r := range rs {
		wg.Add(1)
		go func(r Request) {
			defer wg.Done()
			resp := r.send()

			// Безопасно добавляем ответ в глобальную переменную res
			mu.Lock()
			res.add(resp)
			mu.Unlock()
		}(r)
	}

	wg.Wait()
}

func main() {
	fmt.Println("Hello, World!")

	start()

	addRequest("https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	addRequest("https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	addRequest("https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	addRequest("https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	addRequest("https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	addRequest("https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	addRequest("https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	addRequest("https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	addRequest("https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)

	send()

	// Выводим результаты
	fmt.Printf("Responses count: %d\n", len(res))
	for i, response := range res {
		fmt.Printf("Response %d: Status Code %d, Body: %s\n", i+1, response.statusCode, response.body)
		fmt.Println(response.headers)
	}
}
