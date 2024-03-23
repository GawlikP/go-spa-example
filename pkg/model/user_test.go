package model

import (
  "testing"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "time"
  "os"
)

const timeLayout = time.RFC3339

func TestUserModel(t *testing.T) {
  conn := db.CreateTestConnection(os.Getenv("TEST_MODEL_DB"))
  defer conn.Close()
  db.MigrateUP(conn, "file://../../migrations")

  t.Run("#CreateUser Should add user to database", func(t *testing.T) {
    expectedUser := User{
      // will not be taken by the creation process
      // only use to validate the returned user value
      ID: 1,
      Nickname: "nick",
      Email: "test@mail.com",
      Password: "test",
    }
    user, err := CreateUser(conn, expectedUser)
    if err != nil {
      t.Fatal(err)
    }
    validateUserWithExpected(user, expectedUser, t)
  })
  db.ClearTestDatabase(conn)
}

func validateUserWithExpected(received, expected User, t *testing.T)  {
  if received.ID != expected.ID {
    t.Errorf("User ID should not be returned, got %v want %v", received.ID, expected.ID)
  }
  if received.Nickname != expected.Nickname {
    t.Errorf("User Nickname should not be returned, got %v want %v", received.Nickname, expected.Nickname)
  }
  if received.Email != expected.Email {
    t.Errorf("User Email should not be returned, got %v want %v", received.Email, expected.Email)
  }
  if received.Password != expected.Password {
    t.Errorf("User Password should not be returned, got %v want %v", received.Password, expected.Password)
  }
  if received.CreatedAt != "" {
    _, err := time.Parse(timeLayout, received.CreatedAt)
    if err != nil {
      t.Errorf("User CreatedAt should be set to a valid timestamp, got error %v", err)
    }
  }
  if received.UpdatedAt != "" {
    _, err := time.Parse(timeLayout, received.UpdatedAt)
    if err != nil {
      t.Errorf("User UpdatedAt should be set to a valid timestamp, got error %v", err)
    }
  }
}
