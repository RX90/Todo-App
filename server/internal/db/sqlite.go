package db

import (
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewSQLiteDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "./data/todoapp.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %s", err.Error())
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %s", err.Error())
	}

	if err = applyMigrations(); err != nil && err != migrate.ErrNoChange {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %s", err.Error())
	}

	return db, nil
}

func applyMigrations() error {
	m, err := migrate.New(
		"file://server/migrations",
		"sqlite3://data/todoapp.db")
	if err != nil {
		return err
	}

	return m.Up()
}
