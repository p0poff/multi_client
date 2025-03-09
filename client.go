package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"sync"
	// "unsafe"
)

var rs Requests
var res Responses
var mu sync.Mutex
var ptrStore sync.Map

//export start
func start() {
	rs.clear()
	res.clear()
	// freeAll()
}

// func freeAll() {
// 	ptrStore.Range(func(key, value interface{}) bool {
// 		// Вместо C.free используем обертку
// 		ptr := unsafe.Pointer(key.(*C.char))
// 		// Обертка функции free внутри кода Go
// 		freePtr(ptr)
// 		ptrStore.Delete(key)
// 		return true
// 	})
// }

// Создаем локальную функцию вместо экспортированной C-функции
// func freePtr(ptr unsafe.Pointer) {
// 	C.free(ptr)
// }

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
	// ptrStore.Store(cStr, struct{}{})

	return cStr
}

func main() {
	fmt.Println("Hello, World!")

	start()
}
