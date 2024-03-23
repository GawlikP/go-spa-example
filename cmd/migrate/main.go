package main

import (
  "github.com/GawlikP/go-spa-example/pkg/db"
  "github.com/joho/godotenv"
  "log"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  conn := db.CreateConnection()
  defer conn.Close()
  db.MigrateUP(conn, "file://migrations")
}
