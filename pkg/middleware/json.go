package middleware

import (
  "net/http"
)

func JsonMiddleware() Middleware {
  return func(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("Accept", "application/json")
      next(w, r)
    }
  }
}
