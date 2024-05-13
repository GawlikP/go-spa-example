package middleware

import (
  "github.com/joho/godotenv"
  "log"
  "os"
  // "io"
  "testing"
)

func TestMain(m *testing.M) {
  // log.SetOutput(io.Discard)
  err := godotenv.Load("../../.env")
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  m.Run()
  os.Exit(0)
}
