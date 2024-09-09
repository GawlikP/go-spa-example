package model

import (
  "database/sql"
  "fmt"
  "github.com/GawlikP/go-spa-example/pkg/query"
  "log"
)

type Post struct {
  ID        int    `json:"id"`
  Title     string `json:"title"`
  Content   string `json:"content"`
  UserID    int    `json:"user_id"`
  CreatedAt string `json:"created_at"`
  UpdatedAt string `json:"updated_at"`
  Deleted   bool   `json:"deleted"`
}

type PostWithUser struct {
  ID        int         `json:"id"`
  Title     string      `json:"title"`
  Content   string      `json:"content"`
  UserID    int         `json:"user_id"`
  CreatedAt string      `json:"created_at"`
  UpdatedAt string      `json:"updated_at"`
  Deleted   bool        `json:"deleted"`
  User      SecureUser  `json:"user"`
}

type PostError struct {
  Err error
}

func (e *PostError) Error() string {
  return e.Err.Error()
}

func CreatePost(db *sql.DB, post Post) (Post, error) {
  var err error
  var newPost Post
  log.Print("Validating post")
  err = validatePost(db, post)
  if err != nil {
    log.Print(err)
    return Post{}, err
  }
  log.Print("Creating a new post")
  row := db.QueryRow(query.AddPost, post.Title, post.Content, post.UserID)
  err = row.Scan(&newPost.ID, &newPost.Title, &newPost.Content, &newPost.UserID, &newPost.CreatedAt, &newPost.UpdatedAt, &newPost.Deleted)
  if err != nil {
    log.Print(err)
    log.Print("There was an issue during creating a post")
    return Post{}, err
  }
  return newPost, nil
}

func UpdatePost(db *sql.DB, post Post) (Post, error) {
  var updatedPost Post
  var err error
  log.Print("Validating post")
  err = validatePost(db, post)
  if err != nil {
    log.Print(err)
    return Post{}, err
  }
  log.Print("Validating post ownership")
  err = validatePostOwnership(db, post.ID, post.UserID)
  if err != nil {
    log.Print(err)
    return Post{}, err
  }
  log.Print("Updating a post")
  row := db.QueryRow(query.UpdatePost, post.Title, post.Content, post.UserID, post.ID)
  err = row.Scan(&updatedPost.ID, &updatedPost.Title, &updatedPost.Content, &updatedPost.UserID, &updatedPost.CreatedAt, &updatedPost.UpdatedAt, &updatedPost.Deleted)
  if err != nil {
    log.Print(err)
    log.Print("There was an issue during updating a post")
    return Post{}, err
  }
  return updatedPost, nil
}

func DeletePost(db *sql.DB, postID, userID int) error {
  var err error
  log.Print("Validating post ownership")
  err = validatePostOwnership(db, postID, userID)
  if err != nil {
    log.Print(err)
    return err
  }
  log.Print("Deleting a post")
  _, err = db.Exec(query.DeletePost, postID)
  if err != nil {
    if err == sql.ErrNoRows {
      return &PostError{Err: fmt.Errorf("Post with ID %d does not exist", postID)}
    }
    log.Print(err)
    log.Print("There was an issue during deleting a post")
    return err
  }
  return nil
}

func PostsCount(db *sql.DB) (int, error) {
  var count int
  row := db.QueryRow(query.CountPosts)
  err := row.Scan(&count)
  if err != nil {
    log.Print(err)
    log.Print("There was an issue during counting posts")
    return 0, err
  }
  return count, nil
}

func AllPostsWithUsers(db *sql.DB, page, pageSize int) ([]PostWithUser, error) {
  var posts []PostWithUser
  rows, err := db.Query(query.AllPostsWithUsers, pageSize, (page - 1)*pageSize)
  if err != nil {
    if err == sql.ErrNoRows {
      return []PostWithUser{}, nil
    }
    log.Print(err)
    log.Print("There was an issue during fetching posts")
    return []PostWithUser{}, err
  }
  defer rows.Close()
  for rows.Next() {
    var post PostWithUser
    err := rows.Scan(
      &post.ID,
      &post.Title,
      &post.Content,
      &post.UserID,
      &post.CreatedAt,
      &post.UpdatedAt,
      &post.Deleted,
      &post.User.ID,
      &post.User.Email,
      &post.User.Nickname,
      &post.User.CreatedAt,
      &post.User.UpdatedAt,
    )
    if err != nil {
      log.Print(err)
      log.Print("There was an issue during scanning posts")
      return []PostWithUser{}, err
    }
    posts = append(posts, post)
  }
  return posts, nil
}

func validatePost(db *sql.DB, post Post) error {
  var err error
  if post.Title == "" {
    return &PostError{Err: fmt.Errorf("Title cannot be empty")}
  }
  if post.Content == "" {
    return &PostError{Err: fmt.Errorf("Content cannot be empty")}
  }
  if post.UserID == 0 {
    return &PostError{Err: fmt.Errorf("User ID cannot be empty")}
  }
  _, err = FindUser(db, post.UserID)
  if err != nil {
    return &PostError{Err: fmt.Errorf("User with ID %d does not exist", post.UserID)}
  }
  return nil
}

func FindPostWithUser(db *sql.DB, postID int) (PostWithUser, error) {
  var post PostWithUser
  row := db.QueryRow(query.FindPostByID, postID)
  err := row.Scan(
    &post.ID,
    &post.Title,
    &post.Content,
    &post.UserID,
    &post.CreatedAt,
    &post.UpdatedAt,
    &post.Deleted,
    &post.User.ID,
    &post.User.Email,
    &post.User.Nickname,
    &post.User.CreatedAt,
    &post.User.UpdatedAt,
  )
  if err != nil {
    log.Print(err)
    log.Print("There was an issue during fetching a post")
    if err == sql.ErrNoRows {
      return PostWithUser{}, &PostError{Err: fmt.Errorf("Post with ID %d does not exist", postID)}
    }
    return PostWithUser{}, err
  }
  return post, nil
}

func validatePostOwnership(db *sql.DB, postID, userID int) error {
  var err error
  post, err := FindPostWithUser(db, postID)
  if err != nil {
    return err
  }
  if post.UserID != userID {
    return &PostError{Err: fmt.Errorf("User with ID %d is not an owner of the post with ID %d", userID, postID)}
  }
  return nil
}
