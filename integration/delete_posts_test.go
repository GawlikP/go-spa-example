package integration

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "database/sql"
  "fmt"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/session"
  "github.com/GawlikP/go-spa-example/pkg/router"
)

func TestDeletePosts(t *testing.T) {
  conn := db.CreateTestConnection(os.Getenv("TEST_INTEGRATION_DB"))
  defer conn.Close()
  db.MigrateUP(conn, "file://../migrations")
  defer db.ClearTestDatabase(conn)

  user, err := model.CreateUser(conn, model.User{
    Email: "test@mail.com",
    Password: "password1234",
    Nickname: "testDeletePostUser",
  })
  if err != nil {
    t.Fatal(err)
  }
  post, err := model.CreatePost(conn, model.Post{
    Title: "Test post",
    Content: "Test content",
    UserID: user.ID,
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

  t.Run("DELETE /api/v1/posts/{id} Should return 401 status code when session is invalid", func(t *testing.T) {
    rr := deletePostRequest(conn, t, http.Cookie{}, post.ID)
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
    rr = deletePostRequest(conn, t, cookie, post.ID)
    if rr.Code != http.StatusUnauthorized {
      t.Errorf("Expected status code to be 401, got %d", rr.Code)
    }
  })

  t.Run("DELETE /api/v1/posts/{id} Should return 204 status code", func(t *testing.T) {
    rr := deletePostRequest(conn, t, cookie, post.ID)
    if rr.Code != http.StatusOK {
      t.Errorf("Expected status code to be 200, got %d", rr.Code)
    }
    updatedPost, err := model.FindPostWithUser(conn, post.ID)
    if err != nil {
      t.Fatal(err)
    }
    if updatedPost.Deleted != true {
      t.Errorf("Expected post to be deleted")
    }
  })
}

func deletePostRequest(conn *sql.DB, t *testing.T, cookie http.Cookie, postID int) *httptest.ResponseRecorder {
  r := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/posts/%d", postID), nil)
  r.AddCookie(&cookie)
  rr := httptest.NewRecorder()
  router := router.NewRouter(conn)
  router.ServeHTTP(rr, r)
  return rr
}
