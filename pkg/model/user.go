package model

import (
  "database/sql"
  "log"
  "github.com/GawlikP/go-spa-example/pkg/query"
  "errors"
  "regexp"
  "os"
  "encoding/hex"
  "crypto/sha512"
  "fmt"
)
type User struct {
  ID        int    `json:"id"`
  Email     string `json:"email"`
  Password  string `json:"password"`
  Nickname  string `json:"nickname"`
  CreatedAt string `json:"created_at"`
  UpdatedAt string `json:"updated_at"`
}

type UserError struct {
  Err error
}

func (e *UserError) Error() string {
  return e.Err.Error()
}

const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

func CreateUser(db *sql.DB, user User) (User, error) {
  var err error
  var newUser User
  log.Print("Validating user")
  err = validateUser(db, user, false)
  if err != nil {
    log.Print(err)
    return User{}, err
  }
  log.Print("Creating a new user")
  user.Password = EncryptPassword(user.Password)
  row := db.QueryRow(query.AddUser, user.Email, user.Password, user.Nickname)
  err = row.Scan(&newUser.ID, &newUser.Email, &newUser.Password, &newUser.Nickname, &newUser.CreatedAt, &newUser.UpdatedAt)
  if err != nil {
    log.Print(err)
    log.Print("There was an issue during creating a user")
    return User{}, err
  }
  return newUser, nil
}

func FindUser(db *sql.DB, id int) (User, error) {
  var user User
  row := db.QueryRow(query.FindUser, id)
  err := row.Scan(&user.ID, &user.Email, &user.Nickname, &user.Password, &user.CreatedAt, &user.UpdatedAt)
  if err != nil {
    return User{}, err
  }
  return user, nil
}

func CheckUserPassword(db *sql.DB, id int, password string) error {
  var user User
  row := db.QueryRow(query.FindUser, id)
  err := row.Scan(&user.ID, &user.Email, &user.Nickname, &user.Password, &user.CreatedAt, &user.UpdatedAt)
  if err != nil {
    return err
  }
  if user.Password != EncryptPassword(password) {
    return errors.New("Provided password is not valid")
  }

  return nil
}

func EncryptPassword(password string) string {
  password += os.Getenv("SECRET")
  hasher := sha512.New()
  hasher.Write([]byte(password))
  return hex.EncodeToString(hasher.Sum(nil))
}

func validateUser(db *sql.DB, user User, update bool) error {
  if user.Email == "" {
    return &UserError{ Err: errors.New("Email is required") }
  }
  if user.Password == "" {
    return &UserError{ Err: errors.New("Password is required") }
  }
  if user.Nickname == "" {
    return &UserError{ Err: errors.New("Nickname is required") }
  }
  
  match, err := regexp.MatchString(emailRegex, user.Email)
  if !match || err != nil {
    return &UserError{ Err: errors.New("Email is not valid") }
  }
  if len(user.Password) < 8 {
    return &UserError{ Err: errors.New("Password is to short") }
  }
  if !update {
    foundUser, err := findUserByEmailAndNickname(db, user.Email, user.Nickname)
    if foundUser.ID != 0 {
      return &UserError{ Err: errors.New("User with this Email or Nickname already exists") }
    }
    if err != nil && err != sql.ErrNoRows{
      message := fmt.Sprintf("There was an issue with the database: %v", err.Error())
      return &UserError{ Err: errors.New(message) }
    }
  }
  return nil
}

func findUserByEmailAndNickname(db *sql.DB, email, nickname string) (User, error) {
  var user User
  row := db.QueryRow(query.FindUserByEmailAndNickname, email, nickname)
  err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Nickname, &user.CreatedAt, &user.UpdatedAt)
  if err != nil {
    return User{}, err
  }
  return user, nil
}
