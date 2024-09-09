package handler

import (
  "net/http"
  "encoding/json"
  "github.com/GawlikP/go-spa-example/pkg/response"
  "log"
  "io"
)

func RenderError(w http.ResponseWriter, code int, message string) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  response := response.ErrorResponse{Message: message, Code: code}
  json.NewEncoder(w).Encode(response)
}

func HandleJsonError(w http.ResponseWriter,err error) (http.ResponseWriter) {
  log.Printf("Handling json validation error: %v", err)
  if err == io.EOF {
    log.Printf("Cannot process an empty request")
    RenderError(w, 400, "Invalid empty payload")
    return w
  }
  if err == io.ErrUnexpectedEOF {
    log.Printf("Invalid json body")
    RenderError(w, 400, "Invalid request body")
  }
  log.Printf("Failed to decode request body: %v", err)
  RenderError(w, 400, err.Error())
  return w
}
