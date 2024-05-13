package middleware

import (
  "net/http"
)

func AccessMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers", "content-type, authorization, origin, accept, access-control-allow-origin, access-control-allow-methods, access-control-allow-headers")
    next.ServeHTTP(w, r)
  })
}
