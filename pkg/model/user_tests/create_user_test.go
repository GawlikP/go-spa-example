package model

import (
  "testing"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "time"
  "os"
  "strings"
)

const timeLayout = time.RFC3339

func TestCreateUser(t *testing.T) {
  conn := db.CreateTestConnection(os.Getenv("TEST_MODEL_DB"))
  defer conn.Close()
  db.MigrateUP(conn, "file://../../../migrations")

  t.Run("#CreateUser Should add user to database", func(t *testing.T) {
    expectedUser := model.User{
      ID: 1,
      Nickname: "nick",
      Email: "test@mail.com",
      Password: "test",
    }
    user, err := model.CreateUser(conn, expectedUser)
    if err != nil {
      t.Fatal(err)
    }
    validateUserWithExpected(user, expectedUser, t)
  })

  t.Run("#CreateUser should not create an user with not unique email", func(t *testing.T) {
    invalidUser := model.User{
      Nickname: "nick2",
      Email: "test@mail.com",
      Password: "test",
    }
    user, err := model.CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("invalid not unique uesr email", t)
    }
    validateUserWithExpected(user, model.User{}, t)
  })

  t.Run("#CreateUser should not create an user with not unique nickname", func(t *testing.T){
    var err error
    var user model.User
    invalidUser := model.User{
      Nickname: "nick",
      Email: "test2@mail.com",
      Password: "test",
    }
    user, err = model.CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("invalid not unique user nickname", t)
    }
    validateUserWithExpected(user, model.User{}, t)
  })

  t.Run("#CreateUser should not create an user with invalid email", func(t *testing.T){
    var err error
    var user model.User
    invalidUser := model.User{
      Nickname: "nick2",
      Email: "uniqueButInvalidEmail",
      Password: "test",
    }
    user, err = model.CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("invalid email value", t)
    }
    validateUserWithExpected(user, model.User{}, t)
    invalidUser = model.User{
      Nickname: "nick3",
      Email: "",
      Password: "test",
    }
    user, err = model.CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("invalid email value", t)
    }
    validateUserWithExpected(user, model.User{}, t)
  })

  t.Run("#CreateUser should not create an user with invalid  nickname", func(t *testing.T){
    var err error
    var user model.User
    invalidUser := model.User{
      Nickname: "",
      Email: "test4@mail.com",
      Password: "test",
    }
    user, err = model.CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("empty string nickname", t)
    }
    validateUserWithExpected(user, model.User{}, t)
    invalidUser = model.User{
      Nickname: strings.Repeat("a", 256),
      Email: "test5@mail.com",
      Password: "test",
    }
    user, err = model.CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("nickname longer than 255 characters", t)
    }
    validateUserWithExpected(user, model.User{}, t)
  })

  t.Run("#CreateUser should not create an user with invalid password", func(t *testing.T){
    var err error
    var user model.User
    invalidUser := model.User{
      Nickname: "nick6",
      Email: "test6@mail.com",
      Password: "",
    }
    user, err = model.CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("empty string password", t)
    }
    validateUserWithExpected(user, model.User{}, t)

    invalidUser = model.User{
      Nickname: "nick7",
      Email: "test7@mail.com",
      Password: "abc",
    }
    user, err = model.CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("password shorter than 4 characters", t)
    }
    validateUserWithExpected(user, model.User{}, t)
  })

  db.ClearTestDatabase(conn)
}

func validateUserWithExpected(received, expected model.User, t *testing.T)  {
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

func createUserTestValidationError(message string, t *testing.T) {
  t.Errorf("The CreateUser has not returned any error with %v", message)
}
