package middleware

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "github.com/GawlikP/go-spa-example/pkg/session"
  "github.com/GawlikP/go-spa-example/pkg/model"
)

func AuthorizeMiddlewareTest(t *testing.T) {
  conn := db.CreateTestConnection(os.Getenv("TEST_MIDDLEWARE_DB"))
  defer conn.Close()
  db.MigrateUP(conn, "file://../migrations")
  password := "password1234"
  createdUser, err := model.CreateUser(conn, model.User{
    Email: "user@mail.com",
    Password: password,
    Nickname: "user",
  })
  if err != nil {
    t.Fatal(err)
  }
  encryptedSession, err := session.EncryptSessionWithAES(session.SessionContent{ UserID: createdUser.ID })
  if err != nil {
    t.Fatal(err)
  }
  t.Run("Should return 401 when no session cookie is provided", func(t *testing.T) {
    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
    authorizeHandler := AuthorizeMiddleware(conn)
    handler := Chain(nextHandler, authorizeHandler)
    rw := httptest.NewRecorder()
    r, err := http.NewRequest("GET", "/api/v1/users", nil)
    if err != nil {
      t.Fatal(err)
    }
    handler.ServeHTTP(rw, r)
    if rw.Code != 401 {
      t.Errorf("Expected status code 401, got %d", rw.Code)
    }
  })

  t.Run("Should return 401 when session cookie is invalid", func(t *testing.T) {
    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
    authorizeHandler := AuthorizeMiddleware(conn)
    handler := Chain(nextHandler, authorizeHandler)
    rw := httptest.NewRecorder()
    r, err := http.NewRequest("GET", "/api/v1/users", nil)
    if err != nil {
      t.Fatal(err)
    }
    r.AddCookie(&http.Cookie{
      Name: "session",
      Value: "invalid",
    })
    handler.ServeHTTP(rw, r)
    if rw.Code != 401 {
      t.Errorf("Expected status code 401, got %d", rw.Code)
    }
  })

  t.Run("Should return the handler when session cookie is valid", func(t *testing.T) {
    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      w.WriteHeader(http.StatusOK)
    })
    authorizeHandler := AuthorizeMiddleware(conn)
    handler := Chain(nextHandler, authorizeHandler)
    rw := httptest.NewRecorder()
    r, err := http.NewRequest("GET", "/api/v1/users", nil)
    if err != nil {
      t.Fatal(err)
    }
    r.AddCookie(&http.Cookie{
      Name: "session",
      Value: encryptedSession,
    })
    handler.ServeHTTP(rw, r)
    if rw.Code != http.StatusOK {
      t.Errorf("Expected status code 200, got %d", rw.Code)
    }
  })
  db.ClearTestDatabase(conn)
}
