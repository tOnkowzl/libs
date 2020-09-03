package httpx

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tOnkowzl/libs/contextx"
	"github.com/tOnkowzl/libs/logx"
)

// Request for client do
type Request struct {
	URL    string
	Method string
	Body   interface{}
	Header Header

	HideLogRequest         bool
	HideLogResponse        bool
	UnlimitLogRequestBody  bool
	UnlimitLogResponseBody bool

	fullURL    string
	body       []byte
	marshaller Marshaller
}

func (r *Request) init(ctx context.Context, baseURL string) error {
	r.initFullURL(baseURL)
	r.newMarshaller()
	r.initRequireHeaders(ctx)

	if err := r.marshalBody(); err != nil {
		return err
	}

	return nil
}

func (r *Request) marshalBody() error {
	if r.Body == nil {
		return nil
	}

	if s, ok := r.Body.(string); ok {
		r.body = []byte(s)
		return nil
	}

	b, err := r.marshaller.Marshal(r.Body)
	if err != nil {
		return err
	}

	r.body = b

	return nil
}

func (r *Request) newMarshaller() {
	if strings.ToLower(r.Header[ContentType]) == ApplicationJSON {
		r.marshaller = new(JSON)
		return
	}

	if strings.ToLower(r.Header[ContentType]) == ApplicationXML ||
		strings.ToLower(r.Header[ContentType]) == TextXML {
		r.marshaller = new(XML)
		return
	}

	r.marshaller = new(JSON)
}

func (r *Request) addHeader(key, value string) {
	r.Header[key] = value
}

func (r *Request) initRequireHeaders(ctx context.Context) {
	if r.Header == nil {
		r.Header = Header{}
	}

	if _, ok := r.Header[ContentType]; !ok {
		r.addHeader(ContentType, ApplicationJSON)
	}

	if _, ok := r.Header[HeaderXRequestID]; !ok {
		r.addHeader(HeaderXRequestID, contextx.GetID(ctx))
	}
}

func (r *Request) initFullURL(baseurl string) {
	r.fullURL = baseurl + r.URL
}

func (r *Request) logRequestInfo(ctx context.Context) {
	if r.HideLogRequest {
		return
	}

	var body string
	if r.UnlimitLogRequestBody {
		body = string(r.body)
	} else {
		body = logx.LimitMSG(r.body)
	}

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"method": r.Method,
		"url":    r.fullURL,
		"body":   body,
		"header": r.Header,
	}).Info("client do request information")
}

func (r *Request) logResponseInfo(ctx context.Context, err error, b []byte, duration string, res *http.Response) {
	if r.HideLogResponse {
		return
	}

	if err != nil {
		logx.WithContext(ctx).WithFields(logrus.Fields{
			"url":   r.fullURL,
			"error": err,
		}).Info("client do response information")
		return
	}

	var body string
	if r.UnlimitLogResponseBody {
		body = string(b)
	} else {
		body = logx.LimitMSG(b)
	}

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"duration": duration,
		"status":   res.Status,
		"header":   res.Header,
		"body":     body,
		"url":      r.fullURL,
	}).Info("client do response info")
}
