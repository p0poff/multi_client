package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	opts *opts
}

func NewClient(opts *opts) *Client {
	c := &Client{
		opts: opts,
	}
	return c
}

func (c *Client) sendRequest(req Request) (Response, error) {
	client := &http.Client{
		Timeout: time.Duration(c.opts.Timeout) * time.Second,
	}

	if req.Url == "" {
		return Response{
			Id:         req.Id,
			StatusCode: 0,
			Body:       fmt.Sprintf("Error create request: URL is empty"),
			Headers:    map[string]string{},
		}, fmt.Errorf("URL is empty")
	}

	_, err := url.ParseRequestURI(req.Url)
	if err != nil {
		return Response{
			Id:         req.Id,
			StatusCode: 0,
			Body:       fmt.Sprintf("Error create request: Invalid URL format"),
			Headers:    map[string]string{},
		}, fmt.Errorf("Invalid URL format")
	}

	httpReq, err := http.NewRequest(req.Method, req.Url, nil)
	if err != nil {
		return Response{
			Id:         req.Id,
			StatusCode: 0,
			Body:       fmt.Sprintf("Error create request: %v", err),
			Headers:    map[string]string{},
		}, err
	}

	for key, value := range req.Headers {
		httpReq.Header.Add(key, value)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		// Check for timeout
		if errors.Is(err, context.DeadlineExceeded) {
			return Response{
				Id:         req.Id,
				StatusCode: http.StatusRequestTimeout,
				Body:       "Request timed out",
				Headers:    map[string]string{},
			}, err
		}

		return Response{
			Id:         req.Id,
			StatusCode: 0,
			Body:       fmt.Sprintf("Error sending request: %v", err),
			Headers:    map[string]string{},
		}, err
	}
	defer resp.Body.Close()

	var body []byte
	resp.Body.Read(body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{
			Id:         req.Id,
			StatusCode: resp.StatusCode,
			Body:       fmt.Sprintf("Error reading response body: %v", err),
			Headers:    map[string]string{},
		}, err
	}

	response := Response{
		Id:         req.Id,
		StatusCode: resp.StatusCode,
		Body:       string(respBody),
		Headers:    map[string]string{},
	}

	for key, value := range resp.Header {
		response.Headers[key] = value[0]
	}

	return response, nil
}
