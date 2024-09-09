package integration

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "bytes"
  "database/sql"
  "encoding/json"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/router"
  "github.com/GawlikP/go-spa-example/pkg/session"
  "github.com/GawlikP/go-spa-example/pkg/response"
)

func TestCreatePostsHandler(t *testing.T) {
  conn := db.CreateTestConnection(os.Getenv("TEST_INTEGRATION_DB"))
  defer conn.Close()
  db.MigrateUP(conn, "file://../migrations")
  user, err := model.CreateUser(conn, model.User{
    Email: "test123@mail.com",
    Password: "password1234",
    Nickname: "test123",
  })
  if err != nil {
    t.Fatal(err)
  }
  encryptedSession, err := session.EncryptSessionWithAES(session.SessionContent{ UserID: user.ID })
  if err != nil {
    t.Fatal(err)
  }

  cookie := http.Cookie{
    Name: "session_token",
    Value: encryptedSession,
    HttpOnly: true,
    Secure: true,
    SameSite: http.SameSiteStrictMode,
    Path: "/",
  }

  t.Run("POST /api/v1/posts Should return 401 status code when session is invalid", func(t *testing.T) {
    rr := createPostRequestWithBody(conn, t, http.Cookie{}, []byte{})
    if rr.Code != http.StatusUnauthorized {
      t.Errorf("Expected status code to be 401, got %d", rr.Code)
    }
    cookie := http.Cookie{
      Name: "session_token",
      Value: "invalid",
      HttpOnly: true,
      Secure: true,
      SameSite: http.SameSiteStrictMode,
    }
    rr = createPostRequestWithBody(conn, t, cookie, []byte{})
    if rr.Code != http.StatusUnauthorized {
      t.Errorf("Expected status code to be 401, got %d", rr.Code)
    }
  })

  t.Run("POST /api/v1/posts Should return 201 status code and post data with httpOnly cookie", func(t *testing.T) {
    response := response.SecurePostResponse{}
    payload := []byte(`{"title":"Title","content":"Content"}`)
    count, err := model.PostsCount(conn)
    if err != nil {
      t.Fatal(err)
    }
    postsCount, err := model.PostsCount(conn)
    if err != nil {
      t.Fatal(err)
    }
    rr := createPostRequestWithBody(conn, t, cookie, payload)
    if rr.Code != http.StatusCreated {
      t.Errorf("Expected status code to be 201, got %d", rr.Code)
    }
    err = json.NewDecoder(rr.Body).Decode(&response)
    if err != nil {
      t.Fatal(err)
    }
    if response.Title != "Title" {
      t.Errorf("Expected title to be 'Title', got %s", response.Title)
    }
    if response.Content != "Content" {
      t.Errorf("Expected content to be 'Content', got %s", response.Content)
    }
    if response.UserID != user.ID {
      t.Errorf("Expected user ID to be %d, got %d", user.ID, response.UserID)
    }
    if response.Deleted != false {
      t.Errorf("Expected user deleted to be false, got %t", response.Deleted)
    }
    if ncount, _ := model.PostsCount(conn); ncount != count+1 {
      t.Errorf("Expected posts count to be %d, got %d", count+1, ncount)
    }
    afterCreateCount, err := model.PostsCount(conn)
    if err != nil {
      t.Fatal(err)
    }
    if afterCreateCount != postsCount+1 {
      t.Errorf("Expected posts count to be %d, got %d", postsCount+1, afterCreateCount)
    }
  })

  t.Run("POST /api/v1/posts Should return 400 status code when title is missing", func(t *testing.T) {
    payload := []byte(`{"content":"Content"}`)
    rr := createPostRequestWithBody(conn, t, cookie, payload)
    if rr.Code != http.StatusBadRequest {
      t.Errorf("Expected status code to be 400, got %d", rr.Code)
    }
  })

  t.Run("POST /api/v1/posts Should return 400 status code when content is missing", func(t *testing.T) {
    payload := []byte(`{"title":"Title"}`)
    rr := createPostRequestWithBody(conn, t, cookie, payload)
    if rr.Code != http.StatusBadRequest {
      t.Errorf("Expected status code to be 400, got %d", rr.Code)
    }
  })

  defer db.ClearTestDatabase(conn)
}

func createPostRequestWithBody(conn *sql.DB, t *testing.T, cookie http.Cookie, payload []byte) *httptest.ResponseRecorder {
  rr := httptest.NewRecorder()
  req, err := http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(payload))
  if err != nil {
    t.Fatal(err)
  }
  req.AddCookie(&cookie)
  handler := http.Handler(router.NewRouter(conn))
  handler.ServeHTTP(rr, req)
  return rr
}

