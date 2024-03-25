package router

import (
  "net/http"
  "github.com/GawlikP/go-spa-example/pkg/handler"
  "database/sql"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", handler.HealthHandler)
  mux.HandleFunc("POST /api/v1/register", handler.RegisterHandler(db))
	return mux
}
