package util

import (
  "os"
  "crypto/sha512"
  "encoding/hex"
)

func CreatePasswordHash(password string) string {
  salt := []byte(os.Getenv("TEST_MODEL_DB"))
  pass := []byte(password)
  hasher := sha512.New()
  passwordWithSalt := append(pass, salt...)
  hasher.Write(passwordWithSalt)
  passwordBytes := hasher.Sum(nil)
  return hex.EncodeToString(passwordBytes)
}
