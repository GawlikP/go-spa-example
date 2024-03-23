package model

import (
  "database/sql"
  "log"
  "errors"
  "regexp"
  "github.com/GawlikP/go-spa-example/pkg/query"
  "github.com/GawlikP/go-spa-example/pkg/util"
)

type User struct {
  ID        int    `json:"id"`
  Email     string `json:"email"`
  Password  string `json:"password"`
  Nickname  string `json:"nickname"`
  CreatedAt string `json:"created_at"`
  UpdatedAt string `json:"updated_at"`
}

const emailRegex = `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`

func CreateUser(db *sql.DB, user User) (User, error) {
  var err error
  var newUser User
  log.Print("Validating a new User")
  err = validateUser(db, user)
  if err != nil {
    log.Print("User validation failed!")
    return User{}, err
  }
  log.Print("Creating a new user")
  user.Password = util.CreatePasswordHash(user.Password)
  row := db.QueryRow(query.AddUsersQuerry, user.Email, user.Password, user.Nickname)
  err = row.Scan(&newUser.ID, &newUser.Email, &newUser.Password, &newUser.Nickname, &newUser.CreatedAt, &newUser.UpdatedAt)
  if err != nil {
    log.Print(err)
    log.Print("There was an issue durning creating a user")
    return User{}, err
  }
  return newUser, nil
}

func FindUser(db *sql.DB, id int) (User, error) {
  var err error
  var newUser User
  
  log.Print("Fetching an user")
  row := db.QueryRow(query.FindUserById, id)
  err = row.Scan(&newUser.ID, &newUser.Email, &newUser.Nickname, &newUser.Password, &newUser.CreatedAt, &newUser.UpdatedAt)
  if err != nil {
    log.Print(err)
    log.Print("There was an issue durning finding an user")
    return User{}, err
  }

  return newUser, nil
}

func CheckUserPassword(db *sql.DB, userID int, password string) error {
  var err error
  var user User
  row := db.QueryRow(query.FindUserById, userID)
  err = row.Scan(&user.ID, &user.Email, &user.Nickname, &user.Password, &user.CreatedAt, &user.UpdatedAt)
  if err != nil {
    log.Print(err)
    log.Print("There was an issue durning finding an user")
    return err
  }
  if user.Password == util.CreatePasswordHash(password) {
    return nil
  }
  return errors.New("Provided passowrd does not match the user")
}

func findUserByEmailNick(db *sql.DB, user User) (int, error) {
  var id int
  row := db.QueryRow(query.FindUserByEmailAndNickname, user.Email, user.Nickname)
  err := row.Scan(&id)
  if err != nil {
    if err == sql.ErrNoRows {
      return 0, nil
    }
    return 0, err
  }
  return id, nil
}

func validateUser(db *sql.DB, user User) (error) {
  var foundId int
  var err error
  var match bool
  if user.Email == "" {
    return errors.New("User Email cannot be blank")
  }
  if user.Nickname == "" {
    return errors.New("User Nickname cannot be blank")
  }
  if user.Password == "" {
    return errors.New("User Password caonnot be blank")
  }
  match, err =  regexp.MatchString(emailRegex, user.Email)
  if err != nil || !match {
    return errors.New("User Email value is not a vali Email")
  }
  if len(user.Email) > 255 {
    return errors.New("User Email cannot be longer than 255 characters")
  }

  if len(user.Nickname) > 255 {
    return errors.New("User Nickname cannot be longer than 255 characters")
  }
  if len(user.Password) < 4 {
    return errors.New("User Password needs to be atleast 4 characters long")
  }

  foundId, err = findUserByEmailNick(db, user)
  if err != nil {
    return err
  }
  if foundId != 0 {
    return errors.New("User with that credentials already exists")
  }
  return nil
}
