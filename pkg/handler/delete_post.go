package handler

import (
  "net/http"
  "database/sql"
  "log"
  "strconv"
  "encoding/json"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/response"
)

func DeletePostHandler(db *sql.DB) http.HandlerFunc{
  return func(w http.ResponseWriter, r *http.Request){
    var userID int
    var postID int
    var err error
    stringID := r.PathValue("id")
    postID, err = strconv.Atoi(stringID)
    if err != nil {
      log.Printf("Invalid post ID")
      RenderError(w, http.StatusBadRequest, "Invalid post ID")
      return
    }
    userID, ok := r.Context().Value("user_id").(int)
    if !ok {
      log.Printf("Cannot find user in the context")
      RenderError(w, http.StatusUnauthorized, "Cannot find user in the context")
      return
    }
    err = model.DeletePost(db, postID, userID)
    if err != nil {
      if err == sql.ErrNoRows {
        log.Printf("Post not found")
        RenderError(w, http.StatusNotFound, "Post not found")
        return
      }
      log.Printf("Error during deleting a post: %v", err)
      http.Error(w, "Internal server error", http.StatusInternalServerError)
      return
    }
    w.WriteHeader(http.StatusOK)
    res := response.StandardResponse{
      Message: "Post deleted",
      Code: http.StatusOK,
    }
    json.NewEncoder(w).Encode(res)
  }
}
