package session

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "encoding/json"
  "encoding/hex"
  "os"
  "fmt"
  "io"
)

type SessionContent struct {
  UserID int `json:"user_id"`
}

func SessionToString(session SessionContent) (string, error) {
  bytes, err := json.Marshal(session)
  if err != nil {
    return "", err
  }
  return string(bytes), nil
}

func StringToSession(s string) (SessionContent, error) {
  var session SessionContent
  err := json.Unmarshal([]byte(s), &session)
  if err != nil {
    return session, err
  }
  return session, nil
}

func EncryptSessionWithAES(s SessionContent) (string, error) {
  sstring, err := SessionToString(s)
  secret := os.Getenv("COOKIE_SECRET")

  key, _ := hex.DecodeString(secret)
  plaintext := []byte(sstring)

  block, err := aes.NewCipher(key)
  if err != nil {
    return "", err
  }

  aesGCM, err := cipher.NewGCM(block)
  if err != nil {
    return "", nil
  }

  nonce := make([]byte, aesGCM.NonceSize())
  if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
    return "", err
  }

  ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
  return fmt.Sprintf("%x", ciphertext), nil
}

func DecryptSessionWithAES(s string) (SessionContent, error) {
  keyString := os.Getenv("COOKIE_SECRET")
  enc, err := hex.DecodeString(s)
  if err != nil {
    return SessionContent{}, err
  }
  key, err := hex.DecodeString(keyString)
  if err != nil {
    return SessionContent{}, err
  }
  block, err := aes.NewCipher(key)
  if err != nil {
    return SessionContent{}, err
  }

  aesGCM, err := cipher.NewGCM(block)
  if err != nil {
    return SessionContent{}, err
  }

  nonceSize := aesGCM.NonceSize()
  nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

  plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
  if err != nil {
    return SessionContent{}, err
  }

  return StringToSession(string(plaintext))
}
