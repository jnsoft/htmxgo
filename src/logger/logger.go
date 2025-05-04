package logger

import (
	"log"
	"net/http"
	"os"
	"time"
)

type IsoLogWriter struct{}

func (w *IsoLogWriter) Write(p []byte) (n int, err error) {
	timestamp := time.Now().Format("2006-01-02T15:04:05Z07:00")
	formatted := []byte(timestamp + " " + string(p))
	return os.Stdout.Write(formatted)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(lrw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
	})
}
