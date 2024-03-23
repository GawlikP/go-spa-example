package db

import (
  "database/sql"
  _ "github.com/jackc/pgx/v5/stdlib"
  "log"
  "github.com/golang-migrate/migrate/v4"
  "github.com/golang-migrate/migrate/v4/database/postgres"
  "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateUP(db *sql.DB, dir string) {
  var err error
  err = db.Ping()
  if err != nil {
    log.Printf("Failed to ping the database: %v", err)
  }

  driver, err := postgres.WithInstance(db, &postgres.Config{})
  if err != nil {
    log.Fatalf("Failed to create postgres driver: %v", err)
  }

  source, err := (&file.File{}).Open(dir)
  if err != nil {
    log.Fatalf("Failed to open migration files: %v", err)
  }

  m, err := migrate.NewWithInstance("file", source, "postgres", driver)
  if err != nil {
    log.Fatalf("Failed to create migration instance: %v", err)
  }

  err = m.Up()
  if err != nil && err != migrate.ErrNoChange {
    log.Fatalf("Failed to run migrations: %v", err)
  }
}
