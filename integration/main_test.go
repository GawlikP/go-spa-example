package integration

import (
  "testing"
  "os"
  "github.com/joho/godotenv"
  "log"
  // "io"
)

func TestMain(m *testing.M) {
  // log.SetOutput(io.Discard)
  err := godotenv.Load("../.env")
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  m.Run()
  os.Exit(0)
}
