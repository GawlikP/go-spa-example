package session

import (
  "testing"
)

func TestSession(t *testing.T) {
  t.Run("#SessionToString should return json string of session struct", func(t *testing.T) {
    s := SessionContent{ UserID: 1 }
    expected := "{\"user_id\":1}"
    actual, err := SessionToString(s)
    if err != nil {
      t.Fatal(err)
    }
    if expected != actual {
      t.Errorf("Expected %s, got %s", expected, actual)
    }
  })

  t.Run("#EncryptSessionWithAES should return encrypted session", func(t *testing.T) {
    s := SessionContent{ UserID: 1 }
    encryptedSession, err := EncryptSessionWithAES(s)
    if err != nil {
      t.Fatal(err)
    }
    if encryptedSession == "" {
      t.Error("Expected encrypted session, got empty string")
    }

   t.Run("#DecryptSessionWithAES should return decrypted session", func(t *testing.T) {
     s := SessionContent{ UserID: 1 }
     encryptedSession, err := EncryptSessionWithAES(s)
     if err != nil {
       t.Fatal(err)
     }
     decryptedSession, err := DecryptSessionWithAES(encryptedSession)
     if err != nil {
       t.Fatal(err)
     }
     if decryptedSession.UserID != s.UserID {
       t.Errorf("Expected session UserID to be %d, got %d", s.UserID, decryptedSession.UserID)
     }
   })
  })
}
