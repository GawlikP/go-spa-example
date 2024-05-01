package middleware

import (
  "context"
  "net/http"
  "database/sql"
  "github.com/GawlikP/go-spa-example/pkg/session"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "log"
)

func AuthorizeMiddleware(db *sql.DB) Middleware {
  return func(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      hash, err := r.Cookie("session_token")
      if err != nil {
        log.Printf("Authorization Err: %v", err)
        w.WriteHeader(http.StatusUnauthorized)
        return
      }
      session, err := session.DecryptSessionWithAES(hash.Value)
      if err != nil {
        log.Printf("Authorization Err: %v", err)
        w.WriteHeader(http.StatusUnauthorized)
        return
      }
      if session.UserID == 0 {
        log.Printf("Authorization Err: %v", "User not found")
        w.WriteHeader(http.StatusUnauthorized)
        return
      }
      _, err = model.FindUser(db, session.UserID)
      if err != nil {
        log.Printf("Authorization Err: %v", err)
        w.WriteHeader(http.StatusUnauthorized)
        return
      }
      ctx := r.Context()
      ctx = context.WithValue(ctx, "user_id", session.UserID)
      next(w, r.WithContext(ctx))
    }
  }
}

