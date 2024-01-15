package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (w *LogResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(next http.Handler, logger Logger) http.Handler {
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
		logger.Info(str)
	})
}
