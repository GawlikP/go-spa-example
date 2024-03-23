package db

import (
  "database/sql"
  "fmt"
  _ "github.com/jackc/pgx/v5/stdlib"
  "github.com/GawlikP/go-spa-example/pkg/query"
  "os"
  "log"
)
func CreateConnection() *sql.DB {
  host      := os.Getenv("DB_HOST")
  port      := os.Getenv("DB_PORT")
  user      := os.Getenv("DB_USER")
  password  := os.Getenv("DB_PASS")
  dbname    := os.Getenv("DB_NAME")

  if host == "" || port == "" || user == "" || password == "" || dbname == "" {
    log.Fatal("Please set DB_HOST, DB_PORT, DB_USER, DB_PASS, and DB_NAME environment variables")
  }

  connStr := connectionString(host, port, user, password, dbname)
  return connect(connStr)
}

func CreateTestConnection(database string) *sql.DB {
  host      := os.Getenv("DB_TEST_HOST")
  port      := os.Getenv("DB_TEST_PORT")
  user      := os.Getenv("DB_TEST_USER")
  password  := os.Getenv("DB_TEST_PASS")

  if host == "" || port == "" || user == "" || password == "" {
    log.Fatal("Please set DB_TEST_HOST, DB_TEST_PORT, DB_TEST_USER, DB_TEST_PASS environment variables")
  }
  connStr := connectionString(host, port, user, password, database)
  return connect(connStr)
}

func ClearTestDatabase(db *sql.DB) {
  var err error
  err = db.Ping()
  if err != nil {
    log.Printf("Failed to ping the database: %v", err)
  }
  db.Exec(query.ClearTestsDatabaseQuerry)
}

func connect(connectionString string) *sql.DB {
  db, err := sql.Open("pgx", connectionString)
  if err != nil {
    panic(err)
  }
  return db
}

func connectionString(host, port, user, password, dbname string) string {
  return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}
