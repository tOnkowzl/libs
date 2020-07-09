package httpx

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	BasicAuth  BasicAuth
	HTTPClient *http.Client
	BaseURL    string
}

func NewClient(conf Config) (*Client, error) {
	if conf.BaseURL == "" {
		return nil, errors.New("require base url")
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
		BasicAuth:  conf.BasicAuth,
		HTTPClient: conf.HTTPClient,
		BaseURL:    conf.BaseURL,
	}, nil
}

func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	if err := req.init(ctx, c.BaseURL); err != nil {
		return nil, err
	}

	req.logRequestInfo(ctx)

	httpReq, err := c.makeHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	res, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		req.logResponseInfo(ctx, err, nil, "", nil)
		return nil, err
	}
	duration := time.Since(start).String()

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		req.logResponseInfo(ctx, err, nil, "", nil)
		return nil, err
	}

	req.logResponseInfo(ctx, nil, b, duration, res)

	return &Response{
		Response:   res,
		Marshaller: req.marshaller,
		Body:       b,
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
