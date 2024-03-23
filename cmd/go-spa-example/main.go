package main

import (
  "net/http"
  "github.com/GawlikP/go-spa-example/pkg/router"
)

func main() {
  srv := &http.Server{
		Addr:        "0.0.0.0:8080",
		Handler:     router.NewRouter(),
	}

	srv.ListenAndServe()
}
