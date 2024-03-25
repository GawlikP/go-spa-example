package handler

import (
  "net/http"
  "database/sql"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/response"
  "encoding/json"
)

func RegisterHandler(db *sql.DB) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) { 
    var createdUser model.User
    var newUser model.User
    err := json.NewDecoder(r.Body).Decode(&newUser)
    defer r.Body.Close()
    if err != nil {
      w = HandleJsonError(w, err)
      return
    }
    createdUser, err = model.CreateUser(db, newUser)
    if err != nil {
      re, ok := err.(*model.UserError)
      if ok {
        RenderError(w, http.StatusBadRequest, re.Error())
        return
      }
      RenderError(w, http.StatusInternalServerError, err.Error())
      return
    }
    w.WriteHeader(http.StatusCreated)
    res := response.SecureUserResponse{
      ID: createdUser.ID,
      Email: createdUser.Email,
      Nickname: createdUser.Nickname,
      CreatedAt: createdUser.CreatedAt,
      UpdatedAt: createdUser.UpdatedAt,
    }
    json.NewEncoder(w).Encode(res)
  }
}
