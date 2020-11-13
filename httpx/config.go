package httpx

import (
	"net/http"
	"time"
)

// ClientConfig for service
type ClientConfig struct {
	MaxConns           int
	InsecureSkipVerify bool
	Header             Header
	Timeout            time.Duration
	BaseURL            string
	HTTPClient         *http.Client
}
