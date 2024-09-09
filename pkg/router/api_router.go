package router

import (
  "net/http"
  "github.com/GawlikP/go-spa-example/pkg/handler"
  "github.com/GawlikP/go-spa-example/pkg/middleware"
  "database/sql"
)

func ApiRouter(db *sql.DB) http.Handler {
  mux := http.NewServeMux()
  // /api/v1/
  mux.HandleFunc("OPTIONS /", handler.OptionsHandler)
  // /api/v1/health
  mux.HandleFunc("GET /health", handler.HealthHandler)
  // /api/v1/register
  mux.HandleFunc("POST /register", middleware.Chain(handler.RegisterHandler(db), middleware.JsonMiddleware()))
  // /api/v1/login
  mux.HandleFunc("POST /login", middleware.Chain(handler.LoginHandler(db), middleware.JsonMiddleware()))
  // /api/v1/session/authorize
  mux.HandleFunc("GET /session/authorize", middleware.Chain(handler.AuthorizeHandler, middleware.AuthorizeMiddleware(db), middleware.JsonMiddleware()))
  // /api/v1/posts/{id}
  mux.HandleFunc("PATCH /posts/{id}", middleware.Chain(handler.PatchPostHandler(db), middleware.AuthorizeMiddleware(db), middleware.JsonMiddleware()))
  // /api/v1/posts/{id}
  mux.HandleFunc("DELETE /posts/{id}", middleware.Chain(handler.DeletePostHandler(db), middleware.AuthorizeMiddleware(db), middleware.JsonMiddleware()))
  // // /api/v1/posts
  mux.HandleFunc("GET /posts", middleware.Chain(handler.ListPostsHandler(db), middleware.AuthorizeMiddleware(db), middleware.JsonMiddleware()))
  // /api/v1/posts
  mux.HandleFunc("POST /posts", middleware.Chain(handler.CreatePostHandler(db), middleware.AuthorizeMiddleware(db), middleware.JsonMiddleware()))
  return mux
}
