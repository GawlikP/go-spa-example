package integration

import (
  "testing"
  "database/sql"
  "os"
  "net/http"
  "net/http/httptest"
  "encoding/json"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "github.com/GawlikP/go-spa-example/pkg/router"
  "github.com/GawlikP/go-spa-example/pkg/session"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/response"
)

func TestIntegration(t *testing.T) {
  conn := db.CreateTestConnection(os.Getenv("TEST_INTEGRATION_DB"))
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
  expectedResponse := response.StandardResponse{
    Message: "Success",
    Code: 200,
  }
  t.Run("POST /api/v1/session/authorize should return 200 with valid session", func(t *testing.T) {
    res := response.StandardResponse{}
    cookie := http.Cookie{
      Name: "session_token",
      Value: encryptedSession,
      HttpOnly: true,
      Secure: true,
      SameSite: http.SameSiteStrictMode,
      Path: "/",
    }
    rr := authorizeRequestWithBody(conn, t, cookie)
    if rr.Code != http.StatusOK {
      t.Errorf("Expected status code to be 200, got %d", rr.Code)
    }
    err = json.NewDecoder(rr.Body).Decode(&res)
    if err != nil {
      t.Fatal(err)
    }
    if res.Message != expectedResponse.Message {
      t.Errorf("Expected message to be %s, got %s", expectedResponse.Message, res.Message)
    }
    if res.Code != expectedResponse.Code {
      t.Errorf("Expected code to be %d, got %d", expectedResponse.Code, res.Code)
    }
  })

  t.Run("POST /api/v1/session/authorize should return 401 with no valid session", func(t *testing.T) {
    cookie := http.Cookie{}
    rr := authorizeRequestWithBody(conn, t, cookie)
    if rr.Code != http.StatusUnauthorized {
      t.Errorf("Expected status code to be %v, got %d", http.StatusUnauthorized, rr.Code)
    }
    body := rr.Body.String()
    if body != "" {
      t.Errorf("Expected empty body, got %s", body)
    }
  })

  t.Run("POST /api/v1/session/authorize should return 401 with invalid session value", func(t *testing.T) {
    cookie := http.Cookie{
      Name: "session_token",
      Value: "invalid",
      HttpOnly: true,
      Secure: true,
      SameSite: http.SameSiteStrictMode,
      Path: "/",
    }
    rr := authorizeRequestWithBody(conn, t, cookie)
    if rr.Code != http.StatusUnauthorized {
      t.Errorf("Expected status code to be %v, got %d", http.StatusUnauthorized, rr.Code)
    }
    body := rr.Body.String()
    if body != "" {
      t.Errorf("Expected empty body, got %s", body)
    }
  })

  t.Run("POST /api/v1/session/authorize should return 401 with invalid user id", func(t *testing.T) {
    invalidValue, err := session.EncryptSessionWithAES(session.SessionContent{ UserID: -1 })
    if err != nil {
      t.Fatal(err)
    }
    cookie := http.Cookie{
      Name: "session_token",
      Value: invalidValue,
      HttpOnly: true,
      Secure: true,
      SameSite: http.SameSiteStrictMode,
      Path: "/",
    }
    rr := authorizeRequestWithBody(conn, t, cookie)
    if rr.Code != http.StatusUnauthorized {
      t.Errorf("Expected status code to be %v, got %d", http.StatusUnauthorized, rr.Code)
    }
    body := rr.Body.String()
    if body != "" {
      t.Errorf("Expected empty body, got %s", body)
    }
  })
  db.ClearTestDatabase(conn)
}

func authorizeRequestWithBody(conn *sql.DB, t *testing.T, cookie http.Cookie) *httptest.ResponseRecorder {
  rr := httptest.NewRecorder()
  req, err := http.NewRequest("GET", "/api/v1/session/authorize", nil)
  if err != nil {
    t.Fatal(err)
  }
  req.AddCookie(&cookie)
  handler := http.Handler(router.NewRouter(conn))
  handler.ServeHTTP(rr, req)
  return rr
}
