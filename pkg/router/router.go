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
  // api
  mux.Handle("/api/v1/", http.StripPrefix("/api/v1", ApiRouter(db)))
  return mux
}
