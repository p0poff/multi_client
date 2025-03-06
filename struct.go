package main

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
