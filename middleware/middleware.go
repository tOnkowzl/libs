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
		return c.Path() == "/health"
	}

	UnlimitLogRequestBody  bool
	UnlimitLogResponseBody bool
)

const (
	requestInfoMsg  = "echo request information"
	responseInfoMsg = "echo response information"
)

// Skipper skip middleware
type Skipper func(c echo.Context) bool

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() echo.MiddlewareFunc {
	return RecoverWithConfig(middleware.DefaultRecoverConfig)
}

func RecoverWithConfig(config middleware.RecoverConfig) echo.MiddlewareFunc {
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
						logx.WithContext(c.Request().Context()).Errorf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

// RequestID returns a X-Request-ID middleware.
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if DefaultSkipper(c) {
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
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func LogRequestInfo() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if DefaultSkipper(c) {
				return next(c)
			}

			b := []byte{}
			if c.Request().Body != nil {
				b, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(b))

			var body string
			if UnlimitLogRequestBody {
				body = string(b)
			} else {
				body = logx.LimitMSG(b)
			}

			logx.WithContext(c.Request().Context()).WithFields(logrus.Fields{
				"header": c.Request().Header,
				"body":   body,
			}).Info(requestInfoMsg)

			return next(c)
		}
	}
}

func LogResponseInfo() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if DefaultSkipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()

			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}

			b := resBody.Bytes()

			var body string
			if UnlimitLogResponseBody {
				body = string(b)
			} else {
				body = logx.LimitMSG(b)
			}

			logx.WithContext(c.Request().Context()).WithFields(logrus.Fields{
				"header":    c.Response().Header(),
				"body":      body,
				"method":    req.Method,
				"host":      req.Host,
				"path_uri":  req.RequestURI,
				"remote_ip": c.RealIP(),
				"status":    res.Status,
				"duration":  time.Since(start).String(),
			}).Info(responseInfoMsg)

			return nil
		}
	}
}

func Health() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/health" {
				return c.JSON(http.StatusOK, "ok!")
			}
			return next(c)
		}
	}
}

func JWT(i interface{}) echo.MiddlewareFunc {
	return middleware.JWT(i)
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
