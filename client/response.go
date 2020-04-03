package client

import (
	"net/http"
)

// Response for service do result
type Response struct {
	*http.Response
	Marshaller
	body []byte
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
// output:
//  - ErrResponseNil when Response is nil
//  - ErrResponseNilBody when Response Body is nil
//  - Other error when read body or unmarshal by marshaller
func (r *Response) Unmarshal(v interface{}) error {
	return r.Marshaller.Unmarshal(r.body, v)
}
