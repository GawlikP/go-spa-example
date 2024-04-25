package router

import (
  "net/http"
  "github.com/GawlikP/go-spa-example/pkg/handler"
  "github.com/GawlikP/go-spa-example/ui"
  "io/fs"
  "database/sql"
)

func NewRouter(db *sql.DB) http.Handler {
  mux := http.NewServeMux()

  // ui
  staticFS, _ := fs.Sub(ui.StaticFiles, "dist")
  httpFS := http.FileServer(http.FS(staticFS))

  mux.HandleFunc("/", handler.IndexHandler)
  mux.Handle("/static/", httpFS)

  mux.HandleFunc("GET /api/v1/health", handler.HealthHandler)
  mux.HandleFunc("POST /api/v1/register", handler.RegisterHandler(db))
  mux.HandleFunc("POST /api/v1/login", handler.LoginHandler(db))
  return mux
}
