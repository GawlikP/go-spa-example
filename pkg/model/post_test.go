package model

import (
  "testing"
  "time"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "os"
)

func TestPostModel(t *testing.T) {
  var user User
  conn := db.CreateTestConnection(os.Getenv("TEST_MODEL_DB"))
  defer conn.Close()
  db.MigrateUP(conn, "file://../../migrations")
  user = User{
    Email: "test@mail.com",
    Password: "password1234",
    Nickname: "test",
  }
  user, err := CreateUser(conn, user)

  if err != nil {
    t.Fatal(err)
  }
  expectedPost := Post{
    ID: 1,
    Title: "Title",
    UserID: user.ID,
    Content: "Content",
  }
  t.Run("#PostsCount Should return 0 when there are no posts", func(t *testing.T) {
    count, err := PostsCount(conn)
    if err != nil {
      t.Fatal(err)
    }
    if count != 0 {
      t.Errorf("Expected posts count to be 0, got %d", count)
    }
  })
  t.Run("#CreatePost Should add Post to thye database", func(t *testing.T) {
    post, err := CreatePost(conn, expectedPost)
    if err != nil {
      t.Fatal(err)
    }
    validatePostWithExpected(post, expectedPost, t)
    count, err := PostsCount(conn)
    if err != nil {
      t.Fatal(err)
    }
    if count != 1 {
      t.Errorf("Expected posts count to be 1, got %d", count)
    }
  })
  t.Run("#CreatePost Should return error when user does not exist", func(t *testing.T) {
    expectedPost := Post{
      Title: "Title",
      UserID: -1,
      Content: "Content",
    }
    _, err := CreatePost(conn, expectedPost)
    if err == nil {
      t.Error("Expected error, got nil")
    }
  })
  t.Run("#CreatePost Should validate the post title before adding it to the database", func(t *testing.T) {
    expectedPost := Post{
      Title: "",
      UserID: user.ID,
      Content: "Content",
    }
    _, err := CreatePost(conn, expectedPost)
    if err == nil {
      t.Error("Expected error, got nil")
    }
  })
  t.Run("#CreatePost Should validate the post content before adding it to the database", func(t *testing.T) {
    expectedPost := Post{
      Title: "Title",
      UserID: user.ID,
      Content: "",
    }
    _, err := CreatePost(conn, expectedPost)
    if err == nil {
      t.Error("Expected error, got nil")
    }
  })
  t.Run("#AllPostsWithUsers Should return all posts with users pagginated", func(t *testing.T) {
    page := 1
    pageSize := 2
    expectedPostTwo := Post{
      Title: "Title2",
      UserID: user.ID,
      Content: "Content2",
    }
    _, err := CreatePost(conn, expectedPostTwo)
    if err != nil {
      t.Fatal(err)
    }
    posts, err := AllPostsWithUsers(conn, page, pageSize)
    if err != nil {
      t.Fatal(err)
    }
    if len(posts) != 2 {
      t.Errorf("Expected posts count to be 2, got %d", len(posts))
    }
    if posts[1].User.ID != user.ID {
      t.Errorf("Expected user ID to be %d, got %d", user.ID, posts[0].User.ID)
    }
    if posts[1].User.Email != user.Email {
      t.Errorf("Expected user Email to be %s, got %s", user.Email, posts[0].User.Email)
    }
    if posts[1].User.Nickname != user.Nickname {
      t.Errorf("Expected user Nickname to be %s, got %s", user.Nickname, posts[0].User.Nickname)
    }
    if posts[1].Title != "Title" {
      t.Errorf("Expected post Tittle to be Title, got %s", posts[0].Title)
    }
    if posts[1].Content != "Content" {
      t.Errorf("Expected post Content to be Content, got %s", posts[0].Content)
    }
    if posts[1].Deleted != false {
      t.Errorf("Expected post Deleted to be false, got %t", posts[0].Deleted)
    }
    if posts[0].User.ID != user.ID {
      t.Errorf("Expected user ID to be %d, got %d", user.ID, posts[1].User.ID)
    }
    if posts[0].User.Email != user.Email {
      t.Errorf("Expected user Email to be %s, got %s", user.Email, posts[1].User.Email)
    }
    if posts[0].User.Nickname != user.Nickname {
      t.Errorf("Expected user Nickname to be %s, got %s", user.Nickname, posts[1].User.Nickname)
    }
    if posts[0].Title != "Title2" {
      t.Errorf("Expected post Tittle to be Title2, got %s", posts[1].Title)
    }
    if posts[0].Content != "Content2" {
      t.Errorf("Expected post Content to be Content2, got %s", posts[1].Content)
    }
    if posts[0].Deleted != false {
      t.Errorf("Expected post Deleted to be false, got %t", posts[1].Deleted)
    }
  })

  t.Run("#AllPostsWithUsers Should return empty array when there is no page", func(t *testing.T) {
    page := 99999
    pageSize := 2
    posts, err := AllPostsWithUsers(conn, page, pageSize)
    if err != nil {
      t.Fatal(err)
    }
    if len(posts) != 0 {
      t.Errorf("Expected posts count to be 0, got %d", len(posts))
    }
  });


  db.ClearTestDatabase(conn)
}

func validatePostWithExpected(post, expectedPost Post, t *testing.T) {
  if post.ID != expectedPost.ID {
    t.Errorf("Expected ID to be %d, got %d", expectedPost.ID, post.ID)
  }
  if post.Title != expectedPost.Title {
    t.Errorf("Expected Title to be %s, got %s", expectedPost.Title, post.Title)
  }
  if post.UserID != expectedPost.UserID {
    t.Errorf("Expected UserID to be %d, got %d", expectedPost.UserID, post.UserID)
  }
  if post.Content != expectedPost.Content {
    t.Errorf("Expected Content to be %s, got %s", expectedPost.Content, post.Content)
  }
  if post.CreatedAt != "" {
    _, err := time.Parse(timeLayout, post.CreatedAt)
    if err != nil {
      t.Errorf("Expected CreatedAt to be valid time, got %s", post.CreatedAt)
    }
  }
  if post.UpdatedAt != "" {
    _, err := time.Parse(timeLayout, post.UpdatedAt)
    if err != nil {
      t.Errorf("Expected UpdatedAt to be valid time, got %s", post.UpdatedAt)
    }
  }
  if post.Deleted != false || post.Deleted != expectedPost.Deleted {
    if expectedPost.Deleted {
      t.Errorf("Expected Deleted to be %t, got %t", expectedPost.Deleted, post.Deleted)
    } else {
      t.Errorf("Expected Deleted to be false, got %t", post.Deleted)
    }
  }
}
