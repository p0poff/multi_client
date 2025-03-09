package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	id      int
	url     string
	method  string
	body    string
	headers map[string]string
}

type Response struct {
	Id         int               `json:"id"`
	StatusCode int               `json:"status_code"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
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
			Id:         r.id,
			StatusCode: 0,
			Body:       fmt.Sprintf("Error creating request: %v", err),
			Headers:    map[string]string{},
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
			Id:         r.id,
			StatusCode: 0,
			Body:       fmt.Sprintf("Error sending request: %v", err),
			Headers:    map[string]string{},
		}
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{
			Id:         r.id,
			StatusCode: resp.StatusCode,
			Body:       fmt.Sprintf("Error reading response body: %v", err),
			Headers:    map[string]string{},
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
		Id:         r.id,
		StatusCode: resp.StatusCode,
		Body:       string(respBody),
		Headers:    respHeaders,
	}
}

func (r *Response) get_json() string {
	json, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("Error marshaling response: %v", err)
	}
	return string(json)
}
