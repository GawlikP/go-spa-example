package integration

import (
  "testing"
  "os"
  "net/http"
  "net/http/httptest"
  "github.com/GawlikP/go-spa-example/pkg/router"
  "github.com/GawlikP/go-spa-example/pkg/db"
  "github.com/GawlikP/go-spa-example/pkg/response"
  "github.com/GawlikP/go-spa-example/pkg/model"
  "bytes"
  "encoding/json"
  "database/sql"
)

func TestRegisterHandler(t *testing.T) {
  conn := db.CreateTestConnection(os.Getenv("TEST_INTEGRATION_DB"))
  defer conn.Close()
  db.MigrateUP(conn, "file://../migrations")

  t.Run("POST /api/v1/register Should return 201 with valid body", func(t *testing.T) {
    validRequestBody := model.User{
      Email: "test@mail.com.com",
      Password: "password1234",
      Nickname: "nickname",
    }
    jsonBody, err := json.Marshal(validRequestBody)
    if err != nil {
      t.Fatal(err)
    }
    rr := registerRequestWithBody(jsonBody, conn, t)
    if status := rr.Code; status != http.StatusCreated {
      t.Errorf("Server returned wrong status code: got %v want %v", status, http.StatusCreated)
    }
  })

  t.Run("POST /api/v1/register Should return created User with secure response", func(t *testing.T) {
    var createdUser model.User
    var returnedUser response.SecureUserResponse
    invalidRequestBody := model.User{
      Email: "test2@mail.com.com",
      Password: "password1234",
      Nickname: "nickname2",
    }
    jsonBody, err := json.Marshal(invalidRequestBody) 
    if err != nil {
      t.Fatal(err)
    }
    rr := registerRequestWithBody(jsonBody, conn, t)
    resp := rr.Result()
    defer resp.Body.Close()
    err = json.NewDecoder(resp.Body).Decode(&returnedUser)
    if err != nil {
      t.Errorf("Unable to parse response to User struct: %v", err)
    }
    if createdUser.Password != "" {
      t.Errorf("Register handler returned user with a password!: got %v want %v",
        createdUser.Password, "",
      )
    }
    createdUser, err = model.FindUser(conn, returnedUser.ID)
    if err != nil {
      t.Errorf("Unable to find user in database: %v", err)
    }
    if createdUser.Email != returnedUser.Email {
      t.Errorf("Register handler returned wrong email: got %v want %v",
        createdUser.Email, returnedUser.Email,
      )
    }
    if createdUser.Nickname != returnedUser.Nickname {
      t.Errorf("Register handler returned wrong nickname: got %v want %v",
        createdUser.Nickname, returnedUser.Nickname,
      )
    }
    if createdUser.CreatedAt != returnedUser.CreatedAt {
      t.Errorf("Register handler returned wrong created_at: got %v want %v",
        createdUser.CreatedAt, returnedUser.CreatedAt,
      )
    }
    if createdUser.UpdatedAt != returnedUser.UpdatedAt {
      t.Errorf("Register handler returned wrong updated_at: got %v want %v",
        createdUser.UpdatedAt, returnedUser.UpdatedAt,
      )
    }
  })

  t.Run("POST /api/v1/register should return 'Invalid empty payload' and 400 when request is an empty string", func(t *testing.T) {
    var errorResponse response.ErrorResponse
    invalidRequestBody := []byte("")
    rr := registerRequestWithBody(invalidRequestBody, conn, t)
    if status := rr.Code; status != http.StatusBadRequest {
      t.Errorf("Server returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
    errorResponse = resultToErrorResponse(rr, t)
    validateErrorResponse(errorResponse, response.ErrorResponse{Message: "Invalid empty payload", Code: http.StatusBadRequest}, t)
  })

  t.Run("POST /api/v1/register should return 'Invalid json body' and 400 when request is an empty string", func(t *testing.T) {
    var errorResponse response.ErrorResponse
    invalidRequestBody := []byte("{\"email\" : \"")
    rr := registerRequestWithBody(invalidRequestBody, conn, t)
    if status := rr.Code; status != http.StatusBadRequest {
      t.Errorf("Server returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
    errorResponse = resultToErrorResponse(rr, t)
    validateErrorResponse(errorResponse, response.ErrorResponse{Message: "Invalid request body", Code: http.StatusBadRequest}, t)
  })

  t.Run("POST /api/v1/register should return 'Email is required' with 400 when got empty email", func(t *testing.T) {
    var errorResponse response.ErrorResponse
    invalidRequestBody := model.User{
      Email: "",
      Password: "password1234",
      Nickname: "nickname3",
    }
    jsonBody, err := json.Marshal(invalidRequestBody) 
    if err != nil {
      t.Fatal(err)
    }
    rr := registerRequestWithBody(jsonBody, conn, t)
    if status := rr.Code; status != http.StatusBadRequest {
      t.Errorf("Server returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
    errorResponse = resultToErrorResponse(rr, t)
    validateErrorResponse(errorResponse, response.ErrorResponse{Message: "Email is required", Code: http.StatusBadRequest}, t)
  })

  t.Run("POST /api/v1/register should return 'Password is required' with 400 when got empty password", func(t *testing.T) {
    var errorResponse response.ErrorResponse
    invalidRequestBody := model.User{
      Email: "test3@mail.com",
      Password: "",
      Nickname: "nickname3",
    }
    jsonBody, err := json.Marshal(invalidRequestBody) 
    if err != nil {
      t.Fatal(err)
    }
    rr := registerRequestWithBody(jsonBody, conn, t)
    if status := rr.Code; status != http.StatusBadRequest {
      t.Errorf("Server returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
    errorResponse = resultToErrorResponse(rr, t)
    validateErrorResponse(errorResponse, response.ErrorResponse{Message: "Password is required", Code: http.StatusBadRequest}, t)
  })

  t.Run("POST /api/v1/register should return 'Nickaname is required' with 400 when got empty nickname", func(t *testing.T) {
    var errorResponse response.ErrorResponse
    invalidRequestBody := model.User{
      Email: "test3@mail.com",
      Password: "password1234",
      Nickname: "",
    }
    jsonBody, err := json.Marshal(invalidRequestBody) 
    if err != nil {
      t.Fatal(err)
    }
    rr := registerRequestWithBody(jsonBody, conn, t)
    if status := rr.Code; status != http.StatusBadRequest {
      t.Errorf("Server returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
    errorResponse = resultToErrorResponse(rr, t)
    validateErrorResponse(errorResponse, response.ErrorResponse{Message: "Nickname is required", Code: http.StatusBadRequest}, t)
  })

  t.Run("POST /api/v1/register should return 'Email is invalid' with 400 when got invalid email", func(t *testing.T) {
    var errorResponse response.ErrorResponse
    invalidRequestBody := model.User{
      Email: "testail.com",
      Password: "password1234",
      Nickname: "nickname3",
    }
    jsonBody, err := json.Marshal(invalidRequestBody) 
    if err != nil {
      t.Fatal(err)
    }
    rr := registerRequestWithBody(jsonBody, conn, t)
    if status := rr.Code; status != http.StatusBadRequest {
      t.Errorf("Server returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
    errorResponse = resultToErrorResponse(rr, t)
    validateErrorResponse(errorResponse, response.ErrorResponse{Message: "Email is not valid", Code: http.StatusBadRequest}, t)
  })

  t.Run("POST /api/v1/register should return 'User with this Email or Nickname already exists' with 400 when got not unique parameters", func(t *testing.T) {
    // not unique email
    var errorResponse response.ErrorResponse
    invalidRequestBody := model.User{
      Email: "test@mail.com.com",
      Password: "password1234",
      Nickname: "nickname3",
    }
    jsonBody, err := json.Marshal(invalidRequestBody) 
    if err != nil {
      t.Fatal(err)
    }
    rr := registerRequestWithBody(jsonBody, conn, t)
    if status := rr.Code; status != http.StatusBadRequest {
      t.Errorf("Server returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
    errorResponse = resultToErrorResponse(rr, t)
    validateErrorResponse(errorResponse, response.ErrorResponse{Message: "User with this Email or Nickname already exists", Code: http.StatusBadRequest}, t)

    // not unique nickname
    invalidRequestBody = model.User{
      Email: "test3@mail.com.com",
      Password: "password1234",
      Nickname: "nickname",
    }
    jsonBody, err = json.Marshal(invalidRequestBody) 
    if err != nil {
      t.Fatal(err)
    }
    rr = registerRequestWithBody(jsonBody, conn, t)
    if status := rr.Code; status != http.StatusBadRequest {
      t.Errorf("Server returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
    errorResponse = resultToErrorResponse(rr, t)
    validateErrorResponse(errorResponse, response.ErrorResponse{Message: "User with this Email or Nickname already exists", Code: http.StatusBadRequest}, t)
  })

  t.Run("POST /api/v1/register should return 'Password is to short' with 400 when got password shorter than 8 characters", func(t *testing.T) {
    var errorResponse response.ErrorResponse
    invalidRequestBody := model.User{
      Email: "test4@mail.com",
      Password: "1234567",
      Nickname: "nickname4",
    }
    jsonBody, err := json.Marshal(invalidRequestBody)
    if err != nil {
      t.Fatal(err)
    }
    rr := registerRequestWithBody(jsonBody, conn, t)
    if status := rr.Code; status != http.StatusBadRequest {
      t.Errorf("Server returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
    errorResponse = resultToErrorResponse(rr, t)
    validateErrorResponse(errorResponse, response.ErrorResponse{Message: "Password is to short", Code: http.StatusBadRequest}, t)
    invalidRequestBody = model.User{
      Email: "test5@mail.com",
      Password: "1234567",
      Nickname: "nickname5",
    }
    jsonBody, err = json.Marshal(invalidRequestBody)
    if err != nil {
      t.Fatal(err)
    }
    rr = registerRequestWithBody(jsonBody, conn, t)
    if status := rr.Code; status != http.StatusBadRequest {
      t.Errorf("Server returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
  })
  db.ClearTestDatabase(conn)
}

func resultToErrorResponse(rr *httptest.ResponseRecorder, t *testing.T) response.ErrorResponse {
  var e response.ErrorResponse
  resp := rr.Result()
  defer resp.Body.Close()
  err := json.NewDecoder(resp.Body).Decode(&e)
  if err != nil {
    t.Errorf("Unable to parse response to ErrorResponse struct: %v", err)
  }
  return e
}

func registerRequestWithBody(body []byte, conn *sql.DB, t *testing.T) *httptest.ResponseRecorder {
  req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewReader(body))
  if err != nil {
    t.Fatal(err)
  }
  rr := httptest.NewRecorder()
  handler := http.Handler(router.NewRouter(conn))
  handler.ServeHTTP(rr, req)
  return rr
}

func validateErrorResponse(received response.ErrorResponse, expected response.ErrorResponse, t *testing.T) {
  if received.Message != expected.Message {
    t.Errorf("Register handler returned wrong error message: got %v want %v",
      received.Message, expected.Message,
    )
  }
  if received.Code != expected.Code {
    t.Errorf("Register handler returned wrong error code: got %v want %v",
      received.Code, expected.Code,
    )
  }
}
