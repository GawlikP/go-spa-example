package middleware

import (
  "net/http"
)

func AccessMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS, PATCH")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers", "content-type, authorization, origin, accept, access-control-allow-origin, access-control-allow-methods, access-control-allow-headers, access-control-allow-credentials, credentials")
    next.ServeHTTP(w, r)
  })
}
