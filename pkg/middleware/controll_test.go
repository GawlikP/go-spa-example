package middleware

import (
  "net/http"
  "testing"
)

func AccessMiddlewareTest(t *testing.T) {
  t.Run("Should set Access-Control_Allow-Origin header ", func(t *testing.T) {
    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
    handler := AccessMiddleware(nextHandler)
    rw := NewResponseWriter(nil)
    r, err := http.NewRequest("GET", "/api/v1/users", nil)
    if err != nil {
      t.Fatal(err)
    }
    handler.ServeHTTP(rw, r)
    if rw.Header().Get("Access-Control-Allow-Origin") != "*" {
      t.Error("Access-Control-Allow-Origin header not set to *")
    }
  })

  t.Run("Should set Access-Control-Allow-Methods header ", func(t *testing.T) {
    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
    handler := AccessMiddleware(nextHandler)
    rw := NewResponseWriter(nil)
    r, err := http.NewRequest("GET", "/api/v1/users", nil)
    if err != nil {
      t.Fatal(err)
    }
    handler.ServeHTTP(rw, r)
    if rw.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE" {
      t.Error("Access-Control-Allow-Methods header not set to GET, POST, PUT, DELETE")
    }
  })

  t.Run("Should set Access-Control-Allow-Headers header ", func(t *testing.T) {
    expected := "content-type, authorization, origin, accept, access-control-allow-origin, access-control-allow-methods, access-control-allow-headers"
    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
    handler := AccessMiddleware(nextHandler)
    rw := NewResponseWriter(nil)
    r, err := http.NewRequest("GET", "/api/v1/users", nil)
    if err != nil {
      t.Fatal(err)
    }
    handler.ServeHTTP(rw, r)
    if rw.Header().Get("Access-Control-Allow-Headers") != expected {
      t.Errorf("Access-Control-Allow-Headers header not set to %s", expected)
    }
  })

  t.Run("Should set Access-Control-Allow-Credentials header ", func(t *testing.T) {
    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
    handler := AccessMiddleware(nextHandler)
    rw := NewResponseWriter(nil)
    r, err := http.NewRequest("GET", "/api/v1/users", nil)
    if err != nil {
      t.Fatal(err)
    }
    handler.ServeHTTP(rw, r)
    if rw.Header().Get("Access-Control-Allow-Credentials") != "true" {
      t.Error("Access-Control-Allow-Credentials header not set to true")
    }
  })
}
