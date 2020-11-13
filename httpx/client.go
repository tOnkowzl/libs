package httpx

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

func NewClient(conf ClientConfig) (*Client, error) {
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
	duration := time.Since(start).String()
	if err != nil {
		req.logResponseInfo(ctx, err, nil, duration, nil)
		return nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		req.logResponseInfo(ctx, err, nil, duration, res)
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

	if req.BasicAuth != nil {
		httpReq.SetBasicAuth(req.BasicAuth.Username, req.BasicAuth.Password)
	}

	return httpReq, nil
}
