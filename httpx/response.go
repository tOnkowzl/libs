package httpx

import (
	"net/http"
)

// Response for service do result
type Response struct {
	*http.Response
	Marshaller
	Body []byte
}

// IsOK checking httpStatusCode is < 300
func (r *Response) IsOK() bool {
	return r.StatusCode < 300
}

// IsNotOK checking httpStatusCode more than 299
func (r *Response) IsNotOK() bool {
	return !r.IsOK()
}

// Unmarshal data into v
func (r *Response) Unmarshal(v interface{}) error {
	return r.Marshaller.Unmarshal(r.Body, v)
}
