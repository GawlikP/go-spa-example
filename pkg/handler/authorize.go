package handler


import (
  "net/http"
  "encoding/json"
  "github.com/GawlikP/go-spa-example/pkg/response"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  response := response.StandardResponse{
    Message: "Success",
    Code: 200,
  }
  json.NewEncoder(w).Encode(response)
}
