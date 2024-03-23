package main

import (
  "net/http"
  "github.com/GawlikP/go-spa-example/pkg/router"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "github.com/joho/godotenv"
  "log"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  db := db.CreateConnection()
  defer db.Close()
  err = db.Ping()
  if err != nil {
    log.Printf("Failed to ping the database: %v", err)
  }
  srv := &http.Server{
		Addr:        "0.0.0.0:8080",
		Handler:     router.NewRouter(db),
	}

	srv.ListenAndServe()
}
