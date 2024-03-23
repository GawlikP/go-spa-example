package handler

import (
  "net/http"
  "encoding/json"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  response := map[string]string{
    "message": "HearthBeat",
    "code": "200",
  }
  json.NewEncoder(w).Encode(response)
}
