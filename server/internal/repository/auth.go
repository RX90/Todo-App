package repository

import "github.com/jmoiron/sqlx"

type AuthDB struct {
	db *sqlx.DB
}

func newAuthDB(db *sqlx.DB) *AuthDB {
	return &AuthDB{db: db}
}