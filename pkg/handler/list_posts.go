package handler

import (
  "strconv"
  "net/http"
  "log"
  "database/sql"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/response"
  "encoding/json"
)

func ListPostsHandler(db *sql.DB) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var page string
    var pageSize string
    var posts []model.PostWithUser
    var err error
    params := r.URL.Query()
    page = params.Get("page")
    pageSize = params.Get("pageSize")
    if page == "" {
      page = "1"
    }
    if pageSize == "" {
      pageSize = "10"
    }
    pageInt, err := strconv.Atoi(page)
    if err != nil {
      log.Printf("Cannot convert page to int: %v", err)
      RenderError(w, http.StatusBadRequest, "Invalid page")
      return
    }
    pageSizeInt, err := strconv.Atoi(pageSize)
    if err != nil {
      log.Printf("Cannot convert pageSize to int: %v", err)
      RenderError(w, http.StatusBadRequest, "Invalid pageSize")
      return
    }
    if pageInt < 1 {
      log.Printf("Invalid page number")
      RenderError(w, http.StatusBadRequest, "Invalid page number")
      return
    }
    if pageSizeInt < 1 {
      log.Printf("Invalid pageSize number")
      RenderError(w, http.StatusBadRequest, "Invalid pageSize number")
      return
    }
    posts, err = model.AllPostsWithUsers(db, pageInt, pageSizeInt)
    if err != nil {
      log.Printf("Cannot list posts: %v", err)
      RenderError(w, http.StatusInternalServerError, "Cannot list posts")
      return
    }
    postsCount, err := model.PostsCount(db)
    if err != nil {
      log.Printf("Cannot count posts: %v", err)
      RenderError(w, http.StatusInternalServerError, "Cannot count posts")
      return
    }
    response := response.PostsPageResponse{
      Posts: posts,
      CurrentPage: pageInt,
      PageSize: pageSizeInt,
      TotalCount: postsCount,
      TotalPages: postsCount / pageSizeInt + 1,
    }
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
  }
}
