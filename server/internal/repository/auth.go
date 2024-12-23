package repository

import (
	"fmt"
	"time"

	"github.com/RX90/Todo-App/server/internal/user"
	"github.com/jmoiron/sqlx"
)

type AuthDB struct {
	db *sqlx.DB
}

func newAuthDB(db *sqlx.DB) *AuthDB {
	return &AuthDB{db: db}
}

func (r *AuthDB) isUsernameTaken(username string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE username = $1)", usersTable)

	if err := r.db.QueryRow(query, username).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AuthDB) CreateUser(user user.User) error {
	isTaken, err := r.isUsernameTaken(user.Username)
	if err != nil {
		return err
	}
	if isTaken {
		return fmt.Errorf("username '%s' is already taken", user.Username)
	}

	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3)", usersTable)

	_, err = r.db.Exec(query, user.Name, user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthDB) GetUserId(username, password string) (string, error) {
	var id string
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)

	err := r.db.Get(&id, query, username, password)

	return id, err
}

func (r *AuthDB) NewRefreshToken(token, userId string, expiresAt time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (refresh_token, expires_in) values ($1, $2) RETURNING id", tokensTable)
	row := tx.QueryRow(query, token, expiresAt)

	var tokenId string

	if err := row.Scan(&tokenId); err != nil {
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, token_id) values ($1, $2)", usersTokensTable)
	_, err = tx.Exec(query, userId, tokenId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
