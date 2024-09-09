package integration

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "bytes"
  "database/sql"
  "encoding/json"
  "fmt"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/router"
  "github.com/GawlikP/go-spa-example/pkg/session"
)

func TestPatchPostsHandler(t *testing.T) {
  conn := db.CreateTestConnection(os.Getenv("TEST_INTEGRATION_DB"))
  defer conn.Close()
  db.MigrateUP(conn, "file://../migrations")
  defer db.ClearTestDatabase(conn)

  user, err := model.CreateUser(conn, model.User{
    Email: "test@mail.com",
    Password: "password1234",
    Nickname: "testPatchPostUser",
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

  t.Run("PATCH /api/v1/posts/{id} Should return 401 status code when session is invalid", func(t *testing.T) {
    rr := patchPostRequestWithBody(conn, t, http.Cookie{}, []byte{}, post.ID)
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
    rr = patchPostRequestWithBody(conn, t, cookie, []byte{}, post.ID)
    if rr.Code != http.StatusUnauthorized {
      t.Errorf("Expected status code to be 401, got %d", rr.Code)
    }
  })

  t.Run("PATCH /api/v1/posts/{id} Should return 400 status code when request body is invalid", func(t *testing.T) {
    rr := patchPostRequestWithBody(conn, t, cookie, []byte{}, post.ID)
    if rr.Code != http.StatusBadRequest {
      t.Errorf("Expected status code to be 400, got %d", rr.Code)
    }
  })

  t.Run("PATCH /api/v1/posts/{id} Should return 400 status code when post ID is invalid", func(t *testing.T) {
    validBody, err := json.Marshal(model.Post{
      Title: "Title",
      Content: "Content",
    })
    if err != nil {
      t.Fatal(err)
    }
    rr := patchPostRequestWithBody(conn, t, cookie, validBody, 0)
    if rr.Code != http.StatusBadRequest {
      t.Errorf("Expected status code to be 400, got %d", rr.Code)
    }
  })

  t.Run("PATCH /api/v1/posts/{id} Should return 400 status code when user is not an owner of the post", func(t *testing.T) {
    anotherUser, err := model.CreateUser(conn, model.User{
      Email: "anotheremail@mail.com",
      Password: "password1234",
      Nickname: "anotherUser",
    })
    if err != nil {
      t.Fatal(err)
    }
    encryptedSession, err := session.EncryptSessionWithAES(session.SessionContent{ UserID: anotherUser.ID })
    if err != nil {
      t.Fatal(err)
    }
    anotherUserCookie := http.Cookie{
      Name: "session_token",
      Value: encryptedSession,
      HttpOnly: true,
      Secure: true,
      SameSite: http.SameSiteStrictMode,
      Path: "/",
    }
    validBody, err := json.Marshal(model.Post{
      Title: "Title",
      Content: "Content",
    })
    if err != nil {
      t.Fatal(err)
    }
    rr := patchPostRequestWithBody(conn, t, anotherUserCookie, validBody, post.ID)
    if rr.Code != http.StatusBadRequest {
      t.Errorf("Expected status code to be 400, got %d", rr.Code)
    }
  })

  t.Run("PATCH /api/v1/posts/{id} Should return 201 status code and updated post data", func(t *testing.T) {
    validBody, err := json.Marshal(model.Post{
      Title: "Updated title",
      Content: "Updated content",
    })
    if err != nil {
      t.Fatal(err)
    }
    rr := patchPostRequestWithBody(conn, t, cookie, validBody, post.ID)
    if rr.Code != http.StatusCreated {
      t.Errorf("Expected status code to be 201, got %d", rr.Code)
    }
    response := model.Post{}
    err = json.NewDecoder(rr.Body).Decode(&response)
    if err != nil {
      t.Fatal(err)
    }
    if response.Title != "Updated title" {
      t.Errorf("Expected title to be 'Updated title', got %s", response.Title)
    }
    if response.Content != "Updated content" {
      t.Errorf("Expected content to be 'Updated content', got %s", response.Content)
    }

    updatedPost, err := model.FindPostWithUser(conn, post.ID)
    if err != nil {
      t.Fatal(err)
    }
    if updatedPost.Title != "Updated title" {
      t.Errorf("Expected title to be 'Updated title', got %s", updatedPost.Title)
    }
    if updatedPost.Content != "Updated content" {
      t.Errorf("Expected content to be 'Updated content', got %s", updatedPost.Content)
    }
  })
}

func patchPostRequestWithBody(conn *sql.DB, t *testing.T, cookie http.Cookie, body []byte, postID int) *httptest.ResponseRecorder {
  rr := httptest.NewRecorder()
  url := fmt.Sprintf("/api/v1/posts/%d", postID)
  req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
  if err != nil {
    t.Fatal(err)
  }
  req.AddCookie(&cookie)
  handler := http.Handler(router.NewRouter(conn))
  handler.ServeHTTP(rr, req)
  return rr
}
