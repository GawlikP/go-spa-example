package integration

import (
  "testing"
  "net/http/httptest"
  "os"
  "log"
  "net/http"
  "database/sql"
  "encoding/json"
  "fmt"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/router"
  "github.com/GawlikP/go-spa-example/pkg/session"
  "github.com/GawlikP/go-spa-example/pkg/response"
)

func TestListPostsHandler(t *testing.T) {
  conn := db.CreateTestConnection(os.Getenv("TEST_INTEGRATION_DB"))
  defer conn.Close()
  db.MigrateUP(conn, "file://../migrations")
  defer db.ClearTestDatabase(conn)
  user, err := model.CreateUser(conn, model.User{
    Email: "test1234@mail.com",
    Password: "password1234",
    Nickname: "test1234",
  })
  if err != nil {
    t.Fatal(err)
  }
  postOne, err := model.CreatePost(conn, model.Post{
    Title: "Test post",
    Content: "Test content",
    UserID: user.ID,
  })
  if err != nil {
    t.Fatal(err)
  }
  postTwo, err := model.CreatePost(conn, model.Post{
    Title: "Test post 2",
    Content: "Test content 2",
    UserID: user.ID,
  })
  if err != nil {
    t.Fatal(err)
  }
  postThree, err := model.CreatePost(conn, model.Post{
    Title: "Test post 3",
    Content: "Test content 3",
    UserID: user.ID,
  })
  if err != nil {
    t.Fatal(err)
  }
  postFour, err := model.CreatePost(conn, model.Post{
    Title: "Test post 4",
    Content: "Test content 4",
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
  }

  t.Run("GET /api/v1/posts Should return 401 status code when session is invalid", func(t *testing.T) {
    rr := listPostsRequest(conn, t, http.Cookie{}, "")
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
    rr = listPostsRequest(conn, t, cookie, "")
    if rr.Code != http.StatusUnauthorized {
      t.Errorf("Expected status code to be 401, got %d", rr.Code)
    }
  })

  t.Run("GET /api/v1/posts Should return 200 status code and list of posts", func(t *testing.T) {
    response := response.PostsPageResponse{}
    rr := listPostsRequest(conn, t, cookie, "")
    if rr.Code != http.StatusOK {
      t.Errorf("Expected status code to be 200, got %d", rr.Code)
    }
    err = json.NewDecoder(rr.Body).Decode(&response)
    if err != nil {
      t.Fatal(err)
    }
    log.Printf("Response: %v", response)
    if len(response.Posts) != 4 {
      t.Errorf("Expected 2 posts, got %d", len(response.Posts))
    }
    if response.Posts[3].ID != postOne.ID {
      t.Errorf("Expected post ID to be %d, got %d", postOne.ID, response.Posts[3].ID)
    }
    if response.Posts[2].ID != postTwo.ID {
      t.Errorf("Expected post ID to be %d, got %d", postTwo.ID, response.Posts[2].ID)
    }
    if response.Posts[1].ID != postThree.ID {
      t.Errorf("Expected post ID to be %d, got %d", postThree.ID, response.Posts[1].ID)
    }
    if response.Posts[0].ID != postFour.ID {
      t.Errorf("Expected post ID to be %d, got %d", postFour.ID, response.Posts[0].ID)
    }
    // validate post one with user
    if response.Posts[1].Title != postThree.Title {
      t.Errorf("Expected post title to be %s, got %s", postThree.Title, response.Posts[1].Title)
    }
    if response.Posts[1].Content != postThree.Content {
      t.Errorf("Expected post content to be %s, got %s", postThree.Content, response.Posts[1].Content)
    }
    if response.Posts[1].User.ID != user.ID {
      t.Errorf("Expected user ID to be %d, got %d", user.ID, response.Posts[1].User.ID)
    }
    if response.Posts[1].User.Email != user.Email {
      t.Errorf("Expected user email to be %s, got %s", user.Email, response.Posts[1].User.Email)
    }
    if response.Posts[1].User.Nickname != user.Nickname {
      t.Errorf("Expected user nickname to be %s, got %s", user.Nickname, response.Posts[1].User.Nickname)
    }
    if response.Posts[1].Deleted != false {
      t.Errorf("Expected post deleted to be false, got %t", response.Posts[1].Deleted)
    }
    // validate post two with user
    if response.Posts[0].Title != postFour.Title {
      t.Errorf("Expected post title to be %s, got %s", postFour.Title, response.Posts[0].Title)
    }
    if response.Posts[0].Content != postFour.Content {
      t.Errorf("Expected post content to be %s, got %s", postFour.Content, response.Posts[0].Content)
    }
    if response.Posts[0].User.ID != user.ID {
      t.Errorf("Expected user ID to be %d, got %d", user.ID, response.Posts[0].User.ID)
    }
    if response.Posts[0].User.Email != user.Email {
      t.Errorf("Expected user email to be %s, got %s", user.Email, response.Posts[0].User.Email)
    }
    if response.Posts[0].User.Nickname != user.Nickname {
      t.Errorf("Expected user nickname to be %s, got %s", user.Nickname, response.Posts[0].User.Nickname)
    }
    if response.Posts[0].Deleted != false {
      t.Errorf("Expected post deleted to be false, got %t", response.Posts[0].Deleted)
    }
    // page details
    if response.CurrentPage != 1 {
      t.Errorf("Expected page to be 1, got %d", response.CurrentPage)
    }
    if response.PageSize != 10 {
      t.Errorf("Expected page size to be 10, got %d", response.PageSize)
    }
    if response.TotalPages != 1 {
      t.Errorf("Expected total pages to be 1, got %d", response.TotalPages)
    }
    if response.TotalCount != 4 {
      t.Errorf("Expected total posts to be 2, got %d", response.TotalCount)
    }
  })

  t.Run("GET /api/v1/posts should be able to return valid posts by page", func(t *testing.T) {
  });
}

func listPostsRequest(conn *sql.DB, t *testing.T, cookie http.Cookie, urlSufix string) *httptest.ResponseRecorder {
  rr := httptest.NewRecorder()
  url := fmt.Sprintf("/api/v1/posts%s", urlSufix)
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    t.Fatal(err)
  }
  req.AddCookie(&cookie)
  handler := http.Handler(router.NewRouter(conn))
  handler.ServeHTTP(rr, req)
  return rr
}
