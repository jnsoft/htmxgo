package logger

import (
	"bytes"
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

type extensiveLoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (lrw *extensiveLoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *extensiveLoggingResponseWriter) Write(p []byte) (int, error) {
	lrw.body.Write(p)
	return lrw.ResponseWriter.Write(p)
}

func ExtensiveLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming Request: %s %s %s\nHeaders: %v\n", r.Method, r.URL.Path, r.Proto, r.Header)
		lrw := &extensiveLoggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK, body: &bytes.Buffer{}}
		start := time.Now()
		next.ServeHTTP(lrw, r)
		log.Printf("Outgoing Response: Status: %d\nHeaders: %v\nBody: %s\nDuration: %s\n",
			lrw.statusCode, w.Header(), lrw.body.String(), time.Since(start))
	})
}
