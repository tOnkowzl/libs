package middleware

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/tOnkowzl/libs/contextx"
	"github.com/tOnkowzl/libs/logx"
)

var (
	// DefaultSkipper default of skipper
	DefaultSkipper = func(c echo.Context) bool {
		return c.Path() == "/builds"
	}
)

const (
	requestInfoMsg       = "echo request information"
	responseInfoMsg      = "echo response information"
	logKeywordDontChange = "api_summary"
)

// Skipper skip middleware
type Skipper func(c echo.Context) bool

type Middleware struct {
	Service string
	Skipper Skipper

	UnlimitLogRequestBody  bool
	UnlimitLogResponseBody bool
}

func New(service string) *Middleware {
	return &Middleware{
		Service: service,
		Skipper: DefaultSkipper,
	}
}

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func (m *Middleware) Recover() echo.MiddlewareFunc {
	return m.RecoverWithConfig(middleware.DefaultRecoverConfig)
}

func (m *Middleware) RecoverWithConfig(config middleware.RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = middleware.DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						logx.WithID(m.XRequestID(c)).Errorf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

// Logger log request information
func (m *Middleware) Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if m.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()

			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}

			logx.WithID(m.XRequestID(c)).WithFields(logrus.Fields{
				"method":    req.Method,
				"host":      req.Host,
				"path_uri":  req.RequestURI,
				"remote_ip": c.RealIP(),
				"status":    res.Status,
				"duration":  time.Since(start).String,
				"service":   m.Service,
			}).Info(logKeywordDontChange)

			return
		}
	}
}

// RequestID returns a X-Request-ID middleware.
func (m *Middleware) RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if m.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = uuid.New().String()
			}
			res.Header().Set(echo.HeaderXRequestID, rid)

			ctx := context.WithValue(c.Request().Context(), contextx.ID, rid)
			c.Request().WithContext(ctx)

			return next(c)
		}
	}
}

func (m *Middleware) LogRequestInfo() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if m.Skipper(c) {
				return next(c)
			}

			b := []byte{}
			if c.Request().Body != nil {
				b, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(b))

			var body string
			if m.UnlimitLogRequestBody {
				body = string(b)
			} else {
				body = logx.LimitMSG(b)
			}

			logx.WithID(m.XRequestID(c)).WithFields(logrus.Fields{
				"header": c.Request().Header,
				"body":   body,
			}).Info(requestInfoMsg)

			return next(c)
		}
	}
}

func (m *Middleware) LogResponseInfo() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if m.Skipper(c) {
				return next(c)
			}

			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			if err := next(c); err != nil {
				c.Error(err)
			}

			b := resBody.Bytes()

			var body string
			if m.UnlimitLogResponseBody {
				body = string(b)
			} else {
				body = logx.LimitMSG(b)
			}

			logx.WithID(m.XRequestID(c)).WithFields(logrus.Fields{
				"header": c.Response().Header(),
				"body":   body,
			}).Info(responseInfoMsg)

			return nil
		}
	}
}

func (m *Middleware) Build(buildstamp, githash string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/builds" {
				return c.JSON(http.StatusOK, map[string]string{
					"buildstamp": buildstamp,
					"githash":    githash,
				})
			}
			return next(c)
		}
	}
}

func (m *Middleware) XRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
