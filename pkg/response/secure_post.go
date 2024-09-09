package response

type SecurePostResponse struct {
  ID        int    `json:"id"`
  Title     string `json:"title"`
  Content   string `json:"content"`
  UserID    int    `json:"user_id"`
  CreatedAt string `json:"created_at"`
  UpdatedAt string `json:"updated_at"`
  Deleted   bool   `json:"deleted"`
}
