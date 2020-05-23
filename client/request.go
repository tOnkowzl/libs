package client

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tOnkowzl/libs/logx"
)

// Request for client do
type Request struct {
	URL        string
	Method     string
	XRequestID string
	Body       interface{}
	Header     Header

	HideLogRequest         bool
	HideLogResponse        bool
	UnlimitLogRequestBody  bool
	UnlimitLogResponseBody bool

	fullURL    string
	body       []byte
	marshaller Marshaller
}

func (r *Request) init(baseURL string) {
	r.initXRequestID()
	r.initFullURL(baseURL)
	r.newMarshaller()
	r.initRequireHeaders()
}

func (r *Request) marshalBody() error {
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
	}

	if strings.ToLower(r.Header[ContentType]) == ApplicationXML ||
		strings.ToLower(r.Header[ContentType]) == TextXML {
		r.marshaller = new(XML)
		return
	}

}

func (r *Request) addHeader(key, value string) {
	if r.Header == nil {
		r.Header = Header{}
	}

	r.Header[key] = value
}

func (r *Request) initRequireHeaders() {
	if _, ok := r.Header[ContentType]; !ok {
		r.addHeader(ContentType, ApplicationJSON)
	}

	if _, ok := r.Header[HeaderXRequestID]; !ok {
		r.addHeader(HeaderXRequestID, r.XRequestID)
	}
}

func (r *Request) initXRequestID() {
	if r.XRequestID == "" {
		r.XRequestID = uuid.New().String()
	}
}

func (r *Request) initFullURL(baseurl string) {
	r.fullURL = baseurl + r.URL
}

func (r *Request) logRequestInfo() {
	if r.HideLogRequest {
		return
	}

	var body string
	if r.UnlimitLogRequestBody {
		body = string(r.body)
	} else {
		body = logx.LimitMSG(r.body)
	}

	logx.WithID(r.XRequestID).WithFields(logrus.Fields{
		"method": r.Method,
		"url":    r.fullURL,
		"body":   body,
		"header": r.Header,
	}).Info("client do request info")
}

func (r *Request) logResponseInfo(b []byte, latency string, res *http.Response) {
	if r.HideLogResponse {
		return
	}

	var body string
	if r.UnlimitLogResponseBody {
		body = string(b)
	} else {
		body = logx.LimitMSG(b)
	}

	logx.WithID(r.XRequestID).WithFields(logrus.Fields{
		"latency": latency,
		"status":  res.Status,
		"header":  res.Header,
		"body":    body,
		"url":     r.fullURL,
	}).Info("client do response info")
}
