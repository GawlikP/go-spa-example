package response

type StandardResponse struct {
  Message string `json:"message"`
  Code    int    `json:"code"`
}
