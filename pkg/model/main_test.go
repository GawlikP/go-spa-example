package model

import (
  "testing"
  "os"
  "github.com/joho/godotenv"
  "log"
  "io"
)

func TestMain(m *testing.M) {
  log.SetOutput(io.Discard)
  err := godotenv.Load("../../.env")
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  code := m.Run()
  os.Exit(code)
}
