package router

import (
  "net/http"
  "github.com/GawlikP/go-spa-example/pkg/handler"
  "database/sql"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", handler.HealthHandler)
	return mux
}
