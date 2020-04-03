package client

import (
	"net/http"
	"strings"
	"time"

	"gitdev.inno.ktb/go/libs.git/logx"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Request for client do
type Request struct {
	URL        string
	Method     string
	XRequestID string
	Body       interface{}
	Header     Header

	HideLogRequest      bool
	HideLogResponse     bool
	HideLogRequestBody  bool
	HideLogResponseBody bool

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

func (r *Request) bodyString() string {
	return string(r.body)
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
	if strings.ToLower(r.Header[ContentType]) == ApplicationXML {
		r.marshaller = XML{}
		return
	}
	if strings.ToLower(r.Header[ContentType]) == TextXML {
		r.marshaller = XML{}
		return
	}

	r.marshaller = JSON{}
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
	if !r.HideLogRequest {
		body := r.bodyString()
		if r.HideLogRequestBody {
			body = ""
		}

		logx.WithID(r.XRequestID).WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.fullURL,
			"body":   body,
			"header": r.Header,
		}).Info("client do request info")
	}
}

func (r *Request) logResponseInfo(b []byte, start time.Time, res *http.Response) {
	if !r.HideLogResponse {
		body := string(b)
		if r.HideLogResponseBody {
			body = ""
		}

		logx.WithID(r.XRequestID).WithFields(logrus.Fields{
			"latency": time.Since(start).String(),
			"status":  res.Status,
			"header":  res.Header,
			"body":    body,
			"url":     r.fullURL,
		}).Info("client do response info")
	}
}
