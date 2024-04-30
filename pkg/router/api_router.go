package router

import (
  "net/http"
  "github.com/GawlikP/go-spa-example/pkg/handler"
  "database/sql"
)

func ApiRouter(db *sql.DB) http.Handler {
  mux := http.NewServeMux()
  // /api/v1/
  mux.HandleFunc("OPTIONS /", handler.OptionsHandler)
  // /api/v1/health
  mux.HandleFunc("GET /health", handler.HealthHandler)
  // /api/v1/register
  mux.HandleFunc("POST /register", handler.RegisterHandler(db))
  // /api/v1/login
  mux.HandleFunc("POST /login", handler.LoginHandler(db))
  return mux
}
