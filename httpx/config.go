package httpx

import (
	"net/http"
	"time"
)

// Config for service
type Config struct {
	MaxConns           int
	InsecureSkipVerify bool
	BasicAuth          BasicAuth
	Header             Header
	Timeout            time.Duration
	BaseURL            string
	HTTPClient         *http.Client

	// true is turn on
	TurnONCircuitBreaker   bool
	ErrorPercentThreshold  int
	RequestVolumeThreshold int
	SleepWindow            int
}

// BasicAuth holding username and password for set in *http.Request
type BasicAuth struct {
	Username string
	Password string
}

// HasBasicAuth return config has set username and password for BasicAuth
func (b BasicAuth) HasBasicAuth() bool {
	return b.Username != "" && b.Password != ""
}
