package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	opts *opts
}

func NewServer(opts *opts) *Server {
	s := &Server{
		opts: opts,
	}

	return s
}

func reqFromJson(jsonStr string) (Request, error) {
	var req Request
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		return Request{}, err
	}
	return req, nil

}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.opts.Dbg {
			log.Printf("[DEBUG] Запрос: %s %s", r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func (s *Server) multiSendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var rs Requests
	var wg sync.WaitGroup
	var ress Responses

	err := json.NewDecoder(r.Body).Decode(&rs)
	if err != nil {
		log.Printf("[ERROR] Error decoding JSON: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	for _, r := range rs.Requests {
		wg.Add(1)
		go func(r Request) {
			defer wg.Done()

			c := NewClient(s.opts)

			resp, err := c.sendRequest(r)
			if err != nil {
				log.Printf("[ERROR] Error sending request: %v", err)
			}

			ress.add(resp)
		}(r)
	}

	wg.Wait()

	jsonResponse, err := json.Marshal(ress)
	if err != nil {
		log.Printf("[ERROR] Error marshaling JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))
}

func (s *Server) Start() error {
	addr := ":" + s.opts.Port
	log.Printf("[INFO] Server start! port: %s", addr)

	//routing
	http.Handle("/", s.loggingMiddleware(http.HandlerFunc(s.defaultHandler)))
	http.Handle("/send", s.loggingMiddleware(http.HandlerFunc(s.multiSendHandler)))

	return http.ListenAndServe(addr, nil)
}
