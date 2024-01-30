package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

var locationVerbMap = map[string]string{
	LocationCreate:    http.MethodPost,
	LocationUpdate:    http.MethodPut,
	LocationDelete:    http.MethodDelete,
	LocationListDay:   http.MethodGet,
	LocationListWeek:  http.MethodGet,
	LocationListMonth: http.MethodGet,
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (w *LogResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (mw *Middleware) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLogResponseWriter(w)
		startTime := time.Now()

		next.ServeHTTP(lrw, r)

		str := fmt.Sprintf("%s [%s] %s %s %s %d %s %s",
			r.RemoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.RequestURI(),
			r.Proto,
			lrw.statusCode,
			time.Since(startTime).String(),
			r.UserAgent(),
		)
		mw.logger.Info(str)
	})
}

func (mw *Middleware) requestValidatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		resp := Response{}

		location := path.Base(r.URL.Path)

		if locationVerbMap[location] != r.Method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			resp.Error = http.StatusText(http.StatusMethodNotAllowed)

			err = WriteResponse(w, resp)
			if err != nil {
				mw.logger.Error(err.Error())
			}
			return
		}

		headerContentType := r.Header.Get("Content-Type")

		if !strings.Contains(headerContentType, "application/json") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			resp.Error = http.StatusText(http.StatusUnsupportedMediaType)

			err = WriteResponse(w, resp)
			if err != nil {
				mw.logger.Error(err.Error())
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}

func WriteResponse(w http.ResponseWriter, resp Response) error {
	w.Header().Set("Content-Type", "application/json")

	j, err := json.Marshal(&resp)
	if err != nil {
		return err
	}
	_, err = w.Write(j)
	if err != nil {
		return err
	}

	return nil
}

func MiddlewareChain(handlers ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		for i := len(handlers) - 1; i >= 0; i-- {
			h = handlers[i](h)
		}
		return h
	}
}

type Middleware struct {
	logger Logger
}
