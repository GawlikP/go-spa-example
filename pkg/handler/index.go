package handler

import (
  "net/http"
  "strings"
  "fmt"
  "github.com/GawlikP/go-spa-example/ui"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    w.WriteHeader(http.StatusMethodNotAllowed)
    fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
    return
  }

  if strings.HasPrefix(r.URL.Path, "/api") {
    http.NotFound(w, r)
    return
  }

  if r.URL.Path == "/favicon.ico" {
    rawFile, _ := ui.StaticFiles.ReadFile("dist/favicon.ico")
    w.Write(rawFile)
    return
  }

  rawFile, _ := ui.StaticFiles.ReadFile("dist/index.html")
  w.Write(rawFile)
}
