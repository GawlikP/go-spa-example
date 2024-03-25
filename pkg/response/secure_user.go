package response

type SecureUserResponse struct {
  ID        int    `json:"id"`
  Email     string `json:"email"`
  Nickname  string `json:"nickname"`
  CreatedAt string `json:"created_at"`
  UpdatedAt string `json:"updated_at"`
}
