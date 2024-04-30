package handler

import (
  "net/http"
)

func OptionsHandler(w http.ResponseWriter, r *http.Request) {
  origin := r.Header.Get("Origin")
  w.Header().Set("Access-Control-Allow-Origin", origin)
  w.WriteHeader(http.StatusNoContent)
}
