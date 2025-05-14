package main

type Request struct {
	Id      int               `json:"id"`
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}

type Requests struct {
	Requests []Request `json:"requests"`
}

type Response struct {
	Id         int               `json:"id"`
	StatusCode int               `json:"status_code"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

type Responses struct {
	Responses []Response `json:"responses"`
}

func (r *Responses) add(resp Response) {
	r.Responses = append(r.Responses, resp)
}
