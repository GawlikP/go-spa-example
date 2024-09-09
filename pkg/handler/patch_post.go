package handler

import (
  "net/http"
  "encoding/json"
  "strconv"
  "database/sql"
  "log"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/response"
)

func PatchPostHandler(db *sql.DB) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var userID int
    var postID int
    var post model.Post
    var updatedPost model.Post
    var err error
    err = json.NewDecoder(r.Body).Decode(&post)
    if err != nil {
      http.Error(w, "Invalid request body", http.StatusBadRequest)
      return
    }
    stringID := r.PathValue("id")
    postID, err = strconv.Atoi(stringID)
    if err != nil {
      log.Printf("Invalid post ID")
      RenderError(w, http.StatusBadRequest, "Invalid post ID")
      return
    }
    if err != nil {
      http.Error(w, "Invalid post ID", http.StatusBadRequest)
      return
    }
    userID, ok := r.Context().Value("user_id").(int)
    if !ok {
      log.Printf("Cannot find user in the context")
      RenderError(w, http.StatusUnauthorized, "Cannot find user in the context")
      return
    }
    post.ID = postID
    post.UserID = userID
    updatedPost, err = model.UpdatePost(db, post)
    if err != nil {
      if _, ok := err.(*model.PostError); ok {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }
      log.Printf("Error during updating a post: %v", err)
      http.Error(w, "Internal server error", http.StatusInternalServerError)
      return
    }
    w.WriteHeader(http.StatusCreated)
    res := response.SecurePostResponse{
      ID: updatedPost.ID,
      Title: updatedPost.Title,
      Content: updatedPost.Content,
      CreatedAt: updatedPost.CreatedAt,
      UpdatedAt: updatedPost.UpdatedAt,
      Deleted: updatedPost.Deleted,
    }
    json.NewEncoder(w).Encode(res)
  }
}
