package handler

import (
  "net/http"
  "encoding/json"
  "database/sql"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "github.com/GawlikP/go-spa-example/pkg/response"
  "github.com/GawlikP/go-spa-example/pkg/session"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var foundUser model.User
    var newUser model.User
    err := json.NewDecoder(r.Body).Decode(&newUser)
    defer r.Body.Close()
    if err != nil {
      w = HandleJsonError(w, err)
      return
    }
    err = model.ValidateLoginCredentials(db, newUser)
    if err != nil {
      re, ok := err.(*model.UserError)
      if ok {
        RenderError(w, http.StatusBadRequest, re.Error())
        return
      }
      RenderError(w, http.StatusInternalServerError, err.Error())
      return
    }
    foundUser, err = model.FindUserByEmailAndPassword(db, newUser.Email, newUser.Password)
    if err != nil {
      re, ok := err.(*model.UserError)
      if ok {
        RenderError(w, http.StatusBadRequest, re.Error())
        return
      }
      RenderError(w, http.StatusInternalServerError, err.Error())
      return
    }
    s := session.SessionContent{ UserID: foundUser.ID }
    encryptedSession, err := session.EncryptSessionWithAES(s)
    if err != nil {
      RenderError(w, http.StatusInternalServerError, err.Error())
      return
    }
    http.SetCookie(w, &http.Cookie{
      Name: "session_token",
      Value: encryptedSession,
      HttpOnly: true,
      Secure: true,
      SameSite: http.SameSiteStrictMode,
      Path: "/",
    })

    w.WriteHeader(http.StatusCreated)
    res := response.SecureUserResponse{
      ID: foundUser.ID,
      Email: foundUser.Email,
      Nickname: foundUser.Nickname,
      CreatedAt: foundUser.CreatedAt,
      UpdatedAt: foundUser.UpdatedAt,
    }
    json.NewEncoder(w).Encode(res)
  }
}
