package model

import (
  "testing"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "github.com/GawlikP/go-spa-example/pkg/util"
  "time"
  "os"
  "strings"
)

const timeLayout = time.RFC3339

func TestCreateUser(t *testing.T) {
  conn := db.CreateTestConnection(os.Getenv("TEST_MODEL_DB"))
  defer conn.Close()
  db.MigrateUP(conn, "file://../../migrations")

  t.Run("#CreateUser Should add user to database", func(t *testing.T) {
    expectedUser := User{
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

  t.Run("#CreateUser should not create an user with not unique email", func(t *testing.T) {
    invalidUser := User{
      Nickname: "nick2",
      Email: "test@mail.com",
      Password: "test",
    }
    user, err := CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("invalid not unique uesr email", t)
    }
    validateUserWithExpected(user, User{}, t)
  })

  t.Run("#CreateUser should not create an user with not unique nickname", func(t *testing.T){
    var err error
    var user User
    invalidUser := User{
      Nickname: "nick",
      Email: "test2@mail.com",
      Password: "test",
    }
    user, err = CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("invalid not unique user nickname", t)
    }
    validateUserWithExpected(user, User{}, t)
  })

  t.Run("#CreateUser should not create an user with invalid email", func(t *testing.T){
    var err error
    var user User
    invalidUser := User{
      Nickname: "nick2",
      Email: "uniqueButInvalidEmail",
      Password: "test",
    }
    user, err = CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("invalid email value", t)
    }
    validateUserWithExpected(user, User{}, t)
    invalidUser = User{
      Nickname: "nick3",
      Email: "",
      Password: "test",
    }
    user, err = CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("invalid email value", t)
    }
    validateUserWithExpected(user, User{}, t)
  })

  t.Run("#CreateUser should not create an user with invalid  nickname", func(t *testing.T){
    var err error
    var user User
    invalidUser := User{
      Nickname: "",
      Email: "test4@mail.com",
      Password: "test",
    }
    user, err = CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("empty string nickname", t)
    }
    validateUserWithExpected(user, User{}, t)
    invalidUser = User{
      Nickname: strings.Repeat("a", 256),
      Email: "test5@mail.com",
      Password: "test",
    }
    user, err = CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("nickname longer than 255 characters", t)
    }
    validateUserWithExpected(user, User{}, t)
  })

  t.Run("#CreateUser should not create an user with invalid password", func(t *testing.T){
    var err error
    var user User
    invalidUser := User{
      Nickname: "nick6",
      Email: "test6@mail.com",
      Password: "",
    }
    user, err = CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("empty string password", t)
    }
    validateUserWithExpected(user, User{}, t)

    invalidUser = User{
      Nickname: "nick7",
      Email: "test7@mail.com",
      Password: "abc",
    }
    user, err = CreateUser(conn, invalidUser)
    if err == nil {
      createUserTestValidationError("password shorter than 4 characters", t)
    }
    validateUserWithExpected(user, User{}, t)
  })

  t.Run("#FindUser finds user by provided id", func(t *testing.T) {
    userToFind, err := CreateUser(conn, User{
      Nickname: "nickusertofind",
      Email: "nickusertofind@mail.com",
      Password: "pass",
    })
    if err != nil {
      t.Errorf("Cannot create user to find, err: %v", err)
    }
    foundUser, err := FindUser(conn, userToFind.ID)
    if err != nil {
      t.Errorf("#FindUser method has returned an error while finding a valid user: %v", err)
    }
    // we want a password to match, not the paswords hash
    userToFind.Password = "pass"
    validateUserWithExpected(foundUser, userToFind, t)
  })

  t.Run("#FindUser returns an error when cannnot find user", func(t *testing.T) {
    foundUser, err := FindUser(conn, 0)
    if err == nil {
      t.Error("#FindUser method has not returned an error when looking for not existing user")
    }
    validateUserWithExpected(foundUser, User{}, t)
  })

  t.Run("#CheckUserPassword returns no error when provided password matches users password", func(t *testing.T) {
    password := "password1234"
    userToFind, err := CreateUser(conn, User{
      Nickname: "usercheckpassword",
      Email: "usercheckpassword@mail.com",
      Password: password,
    })
    if err != nil {
      t.Errorf("Cannot create user to find, err: %v", err)
    }
    err = CheckUserPassword(conn, userToFind.ID, password)
    if err != nil {
      t.Errorf("#CheckUserPassword method has returned an error when checking the valid password, err: %v", err)
    }
  })

  t.Run("#CheckUserPassword returns an error when provided password does not match the users password", func(t *testing.T){
    password := "pass098765"
    userToFind, err := CreateUser(conn, User{
      Nickname: "invaliduserpassword",
      Email: "invaliduserpassword@mail.com",
      Password: "pass",
    })
    if err != nil {
      t.Errorf("Cannot create user to find, err: %v", err)
    }
    err = CheckUserPassword(conn, userToFind.ID, password)
    if err == nil {
      t.Error("#CheckUserPassword method has returned no error when checking the invalid password")
    }
  })
  db.ClearTestDatabase(conn)
}

func validateUserWithExpected(received, expected User, t *testing.T)  {
  if received.ID != expected.ID {
    t.Errorf("User IDs does not match, got %v want %v", received.ID, expected.ID)
  }
  if received.Nickname != expected.Nickname {
    t.Errorf("User Nicknames does not match, got %v want %v", received.Nickname, expected.Nickname)
  }
  if received.Email != expected.Email {
    t.Errorf("User Emails does not match, got %v want %v", received.Email, expected.Email)
  }
  if expected.Password != "" {
    t.Log("Password Profided")
    if received.Password != util.CreatePasswordHash(expected.Password) {
      t.Errorf("User Passwords does not match, got %v want %v", received.Password, expected.Password)
    }
  } else {
    t.Log("Password not provided")
    if received.Password != expected.Password {
      t.Errorf("User Passwords does not match, got %v want %v", received.Password, expected.Password)
    }
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

