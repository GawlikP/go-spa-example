package db

import (
  "database/sql"
  "fmt"
  _ "github.com/jackc/pgx/v5/stdlib"
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
    log.Fatal("Please set PG_HOST, PG_PORT, PG_USER, PG_PASSWORD, and PG_DB environment variables")
  }

  connStr := connectionString(host, port, user, password, dbname)
  db, err := sql.Open("pgx", connStr)
  if err != nil {
    panic(err)
  }
  return db
}

func connectionString(host, port, user, password, dbname string) string {
  return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}
