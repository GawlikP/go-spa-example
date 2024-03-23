package router

import (
  "net/http"
  "github.com/GawlikP/go-spa-example/pkg/handler"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", handler.HealthHandler)
	return mux
}
