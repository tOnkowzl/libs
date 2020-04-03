package client

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/google/uuid"
)

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

// Header wrap map[string]string for header
type Header map[string]string

type Clienter interface {
	Do(*Request) (*Response, error)
}

type Client struct {
	BasicAuth             BasicAuth
	HTTPClient            *http.Client
	BaseURL               string
	TurnONCircuitBreaker  bool
	ErrorPercentThreshold int
	HystrixName           string
}

func New(conf Config) (*Client, error) {
	if conf.BaseURL == "" {
		return nil, errors.New("require base url")
	}

	hystrixName := uuid.New().String()
	if conf.TurnONCircuitBreaker {
		hystrix.ConfigureCommand(hystrixName, hystrix.CommandConfig{
			Timeout:                int(conf.Timeout / time.Millisecond),
			MaxConcurrentRequests:  conf.MaxConns * 2,
			ErrorPercentThreshold:  conf.ErrorPercentThreshold,
			RequestVolumeThreshold: conf.RequestVolumeThreshold,
			SleepWindow:            conf.SleepWindow,
		})
	}

	if conf.HTTPClient == nil {
		conf.HTTPClient = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: conf.MaxConns,
				MaxConnsPerHost:     conf.MaxConns,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: conf.InsecureSkipVerify},
			},
			Timeout: conf.Timeout,
		}
	}

	return &Client{
		BasicAuth:             conf.BasicAuth,
		HTTPClient:            conf.HTTPClient,
		BaseURL:               conf.BaseURL,
		TurnONCircuitBreaker:  conf.TurnONCircuitBreaker,
		ErrorPercentThreshold: conf.ErrorPercentThreshold,
		HystrixName:           hystrixName,
	}, nil
}

func (c *Client) Do(req *Request) (*Response, error) {
	req.init(c.BaseURL)

	if err := req.marshalBody(); err != nil {
		return nil, err
	}

	req.logRequestInfo()

	httpReq, err := c.makeHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	res, err := c.do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	req.logResponseInfo(b, start, res)

	return &Response{
		Response:   res,
		Marshaller: req.marshaller,
		body:       b,
	}, nil
}

func (c *Client) makeHTTPRequest(req *Request) (*http.Request, error) {
	httpReq, err := http.NewRequest(req.Method, req.fullURL, bytes.NewReader(req.body))
	if err != nil {
		return nil, err
	}

	for k, v := range req.Header {
		httpReq.Header.Set(k, v)
	}

	if c.BasicAuth.HasBasicAuth() {
		httpReq.SetBasicAuth(c.BasicAuth.Username, c.BasicAuth.Password)
	}

	return httpReq, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	if !c.TurnONCircuitBreaker {
		return c.HTTPClient.Do(req)
	}

	var httpRes *http.Response
	err := hystrix.Do(c.HystrixName, func() error {
		res, err := c.HTTPClient.Do(req)
		if err != nil {
			return err
		}

		httpRes = res

		return nil
	}, nil)
	if err != nil {
		return nil, err
	}

	return httpRes, nil
}
