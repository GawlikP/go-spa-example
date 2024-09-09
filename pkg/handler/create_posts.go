package handler

import (
  "net/http"
  "log"
  "database/sql"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/response"
  "encoding/json"
)

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var userID int
    var newPost model.Post
    err := json.NewDecoder(r.Body).Decode(&newPost)
    defer r.Body.Close()
    if err != nil {
      w = HandleJsonError(w, err)
      return
    }
    userID, ok := r.Context().Value("user_id").(int)
    if ok != true {
      log.Printf("Cannot find user in the context")
      RenderError(w, http.StatusUnauthorized, "Cannot find user in the context")
      return
    }
    newPost.UserID = userID
    createdPost, err := model.CreatePost(db, newPost)
    if err != nil {
      RenderError(w, http.StatusBadRequest, err.Error())
      return
    }
    w.WriteHeader(http.StatusCreated)
    res := response.SecurePostResponse{
      ID: createdPost.ID,
      Title: createdPost.Title,
      UserID: createdPost.UserID,
      Content: createdPost.Content,
      CreatedAt: createdPost.CreatedAt,
      UpdatedAt: createdPost.UpdatedAt,
      Deleted: createdPost.Deleted,
    }
    json.NewEncoder(w).Encode(res)
  }
}
