package model

import (
  "database/sql"
  "log"
  "github.com/GawlikP/go-spa-example/pkg/query"
)

type User struct {
  ID        int    `json:"id"`
  Email     string `json:"email"`
  Password  string `json:"password"`
  Nickname  string `json:"nickname"`
  CreatedAt string `json:"created_at"`
  UpdatedAt string `json:"updated_at"`
}

func CreateUser(db *sql.DB, user User) (User, error) {
	var err error
  var newUser User
	log.Print("Creating a new user")
	row := db.QueryRow(query.AddUsersQuerry, user.Email, user.Password, user.Nickname)
	err = row.Scan(&newUser.ID, &newUser.Email, &newUser.Password, &newUser.Nickname, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		log.Print(err)
		log.Print("There was an issue durning creating a user")
		return User{}, err
	}
	return newUser, nil
}
