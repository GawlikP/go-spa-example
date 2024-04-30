package middleware

import (
  "net/http"
  "log"
  "time"
)

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
    return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rw := NewResponseWriter(w)

        start := time.Now()

        defer func() {
          duration := time.Since(start)
          log.Printf("%v %s: %d %v   \n", r.Method, r.URL.Path, rw.statusCode, duration)
        }()
        next.ServeHTTP(rw, r)
    })
}
