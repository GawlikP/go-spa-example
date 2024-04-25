package integration

import (
  "testing"
  "net/http"
  "net/http/httptest"
  "database/sql"
  "os"
  "strings"
  "bytes"
  "encoding/json"
  "github.com/GawlikP/go-spa-example/pkg/router"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/response"
  "github.com/GawlikP/go-spa-example/pkg/db"
)

func TestLoginHandler(t *testing.T) {
  conn := db.CreateTestConnection(os.Getenv("TEST_INTEGRATION_DB"))
  defer conn.Close()
  db.MigrateUP(conn, "file://../migrations")
  password := "password"
  createdUser, err := model.CreateUser(conn, model.User{
    Email: "user@mail.com",
    Password: password,
    Nickname: "user",
  })
  if err != nil {
    t.Fatal(err)
  }
  t.Run("POST /api/v1/login Should return 201 status code and user data with httpOnly cookie", func(t *testing.T) {
    var res response.SecureUserResponse
    validRequestBody := model.User{
      Email: createdUser.Email,
      Password: password,
    }
    jsonBody, err := json.Marshal(validRequestBody)
    if err != nil {
      t.Fatal(err)
    }
    rr := loginRequestWithBody(jsonBody, conn, t)
    if rr.Code != http.StatusCreated {
      t.Errorf("Expected status code to be 201, got %d", rr.Code)
    }
    if !strings.Contains(rr.Header().Get("Set-Cookie"), "session_token") {
      t.Error("Expected session_token cookie, got none")
    }
    sessionCookie := rr.Header().Get("Set-Cookie")
    if sessionCookie == "" {
      t.Error("Expected session_token cookie, got none")
    }
    if !strings.Contains(sessionCookie, "HttpOnly") {
      t.Error("Expected HttpOnly flag in cookie, got none")
    }
    if !strings.Contains(sessionCookie, "Secure") {
      t.Error("Expected Secure flag in cookie, got none")
    }
    if !strings.Contains(sessionCookie, "SameSite=Strict") {
      t.Error("Expected SameSite=Strict flag in cookie, got none")
    }
    if !strings.Contains(sessionCookie, "Path=/") {
      t.Error("Expected Path=/ flag in cookie, got none")
    }
    err = json.NewDecoder(rr.Body).Decode(&res)
    if err != nil {
      t.Fatal(err)
    }
    if res.ID != createdUser.ID {
      t.Errorf("Expected user ID to be %d, got %d", createdUser.ID, res.ID)
    }
    if res.Email != createdUser.Email {
      t.Errorf("Expected user email to be %s, got %s", createdUser.Email, res.Email)
    }
    if res.Nickname != createdUser.Nickname {
      t.Errorf("Expected user nickname to be %s, got %s", createdUser.Nickname, res.Nickname)
    }
  })

  t.Run("POST /api/v1/login Should return 400 status code when user does not exist", func(t *testing.T) {
    expectedMessage := "Provided credentials are invalid"
    var res response.ErrorResponse
    invalidRequestBody := model.User{
      Email: "invalidUserEmail@mail.com",
      Password: password,
    }
    jsonBody, err := json.Marshal(invalidRequestBody)
    if err != nil {
      t.Fatal(err)
    }
    rr := loginRequestWithBody(jsonBody, conn, t)
    if rr.Code != http.StatusBadRequest {
      t.Errorf("Expected status code to be 400, got %d", rr.Code)
    }
    err = json.NewDecoder(rr.Body).Decode(&res)
    if err != nil {
      t.Fatal(err)
    }
    if res.Message != expectedMessage {
      t.Errorf("Expected error message to be %s not found, got %s", expectedMessage, res.Message)
    }
  })

  t.Run("POST /api/v1/login Should return 400 status code when password is invalid", func(t *testing.T) {
    expectedMessage := "Provided password is not valid"
    var res response.ErrorResponse
    invalidRequestBody := model.User{
      Email: createdUser.Email,
      Password: "pswd",
    }
    jsonBody, err := json.Marshal(invalidRequestBody)
    if err != nil {
      t.Fatal(err)
    }
    rr := loginRequestWithBody(jsonBody, conn, t)
    if rr.Code != http.StatusBadRequest {
      t.Errorf("Expected status code to be 400, got %d", rr.Code)
    }
    err = json.NewDecoder(rr.Body).Decode(&res)
    if err != nil {
      t.Fatal(err)
    }
    if res.Message != expectedMessage {
      t.Errorf("Expected error message to be %s, got %s", expectedMessage, res.Message)
    }
  })

  t.Run("POST /api/v1/login Should return 400 status code when email is not present", func(t *testing.T) {
    expectedMessage := "Email is required"
    var res response.ErrorResponse
    invalidRequestBody := model.User{
      Password: password,
    }
    jsonBody, err := json.Marshal(invalidRequestBody)
    if err != nil {
      t.Fatal(err)
    }
    rr := loginRequestWithBody(jsonBody, conn, t)
    if rr.Code != http.StatusBadRequest {
      t.Errorf("Expected status code to be 400, got %d", rr.Code)
    }
    err = json.NewDecoder(rr.Body).Decode(&res)
    if err != nil {
      t.Fatal(err)
    }
    if res.Message != expectedMessage {
      t.Errorf("Expected error message to be %s, got %s", expectedMessage, res.Message)
    }
  })

  t.Run("POST /api/v1/login Should return 400 status code when password is not present", func(t *testing.T) {
    expectedMessage := "Password is required"
    var res response.ErrorResponse
    invalidRequestBody := model.User{
      Email: createdUser.Email,
    }
    jsonBody, err := json.Marshal(invalidRequestBody)
    if err != nil {
      t.Fatal(err)
    }
    rr := loginRequestWithBody(jsonBody, conn, t)
    if rr.Code != http.StatusBadRequest {
      t.Errorf("Expected status code to be 400, got %d", rr.Code)
    }
    err = json.NewDecoder(rr.Body).Decode(&res)
    if err != nil {
      t.Fatal(err)
    }
    if res.Message != expectedMessage {
      t.Errorf("Expected error message to be %s, got %s", expectedMessage, res.Message)
    }
  })
  db.ClearTestDatabase(conn)
}

func loginRequestWithBody(body []byte, conn *sql.DB, t *testing.T) *httptest.ResponseRecorder {
  req, err := http.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
  if err != nil {
    t.Fatal(err)
  }
  rr := httptest.NewRecorder()
  handler := http.Handler(router.NewRouter(conn))
  handler.ServeHTTP(rr, req)
  return rr
}
