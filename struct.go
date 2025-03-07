package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	url     string
	method  string
	body    string
	headers map[string]string
}

type Response struct {
	statusCode int
	body       string
	headers    map[string]string
}

type Requests []Request
type Responses []Response

func (r *Requests) add(r2 Request) {
	*r = append(*r, r2)
}

func (r *Requests) clear() {
	*r = []Request{}
}

func (r *Responses) add(r2 Response) {
	*r = append(*r, r2)
}

func (r *Responses) clear() {
	*r = []Response{}
}

func (r *Request) send() Response {
	// Создаем HTTP клиент
	client := &http.Client{
		Timeout: time.Second * 10, // Устанавливаем таймаут
	}

	// Создаем запрос
	req, err := http.NewRequest(r.method, r.url, strings.NewReader(r.body))
	if err != nil {
		return Response{
			statusCode: 0,
			body:       fmt.Sprintf("Error creating request: %v", err),
			headers:    map[string]string{},
		}
	}

	// Добавляем заголовки
	for key, value := range r.headers {
		req.Header.Add(key, value)
	}

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return Response{
			statusCode: 0,
			body:       fmt.Sprintf("Error sending request: %v", err),
			headers:    map[string]string{},
		}
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{
			statusCode: resp.StatusCode,
			body:       fmt.Sprintf("Error reading response body: %v", err),
			headers:    map[string]string{},
		}
	}

	// Преобразуем заголовки ответа
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0]
		}
	}

	// Возвращаем Response
	return Response{
		statusCode: resp.StatusCode,
		body:       string(respBody),
		headers:    respHeaders,
	}
}
