package response

import "github.com/GawlikP/go-spa-example/pkg/model"

type PostsPageResponse struct {
  Posts []model.PostWithUser `json:"posts"`
  CurrentPage int `json:"current_page"`
  PageSize int `json:"page_size"`
  TotalPages int `json:"total_pages"`
  TotalCount int `json:"total_count"`
}
