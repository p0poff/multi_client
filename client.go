package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"sync"
	"unsafe"
)

var rs Requests
var res Responses
var mu sync.Mutex
var ptrStore sync.Map

//export start
func start() {
	rs.clear()
	res.clear()
	freeAll()
}

func freeAll() {
	ptrStore.Range(func(key, value interface{}) bool {
		C.free(unsafe.Pointer(key.(*C.char))) // Освобождаем память
		ptrStore.Delete(key)                  // Удаляем из карты
		return true
	})
}

//export addRequest
func addRequest(id int, url *C.char, method *C.char, body *C.char, headers_json *C.char) bool {
	headers := map[string]string{}
	err := json.Unmarshal([]byte(C.GoString(headers_json)), &headers)

	if err != nil {
		return false
	}

	r := Request{
		id:      id,
		url:     C.GoString(url),
		method:  C.GoString(method),
		body:    C.GoString(body),
		headers: headers,
	}
	rs.add(r)

	return true
}

//export send
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

//export getResponsesJson
func getResponsesJson() *C.char {
	cStr := C.CString(res.get_json())
	ptrStore.Store(cStr, struct{}{})

	return cStr
}

func main() {
	fmt.Println("Hello, World!")

	start()

	// addRequest(1, "https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	// addRequest(2, "https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	// addRequest(3, "https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	// addRequest(4, "https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	// addRequest(5, "https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	// addRequest(6, "https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	// addRequest(7, "https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	// addRequest(8, "https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	// addRequest(9, "https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)
	// addRequest(0, "https://httpbin.org/delay/2", "GET", "{}", `{"Content-Type": "application/json"}`)

	send()

	// Выводим результаты
	fmt.Printf("Responses count: %d\n", len(res))

	fmt.Printf("Responses: %s\n", getResponsesJson())

}
