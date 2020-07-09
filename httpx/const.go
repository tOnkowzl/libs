package httpx

import "errors"

var (
	// ErrResponseNil return when unmarshal response but response is nil
	ErrResponseNil = errors.New("response is nil")

	// ErrResponseNilBody return when unmarshal response but response body is nil
	ErrResponseNilBody = errors.New("response body is nil")

	// ErrServiceConfigNotSpecifyContentType return when call service without Content-Type header
	ErrServiceConfigNotSpecifyContentType = errors.New("service header not set Content-Type")

	// ErrContentTypeNotSupport return when call service with not support Content-Type
	ErrContentTypeNotSupport = errors.New("Content-Type not support")
)
