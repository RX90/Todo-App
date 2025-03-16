package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/RX90/Todo-App/server/internal/todo"
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

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE LOWER(username) = LOWER($1))", usersTable)
	err := r.db.QueryRow(query, username).Scan(&exists)

	return exists, err
}

func (r *AuthDB) CreateUser(user todo.User) error {
	isTaken, err := r.isUsernameTaken(user.Username)
	if err != nil {
		return err
	}
	if isTaken {
		return fmt.Errorf("username is already taken")
	}

	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) values ($1, $2)", usersTable)

	_, err = r.db.Exec(query, user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthDB) GetUserId(user todo.User) (string, error) {
	var id string

	query := fmt.Sprintf("SELECT id FROM %s WHERE username = $1 AND password_hash = $2", usersTable)
	err := r.db.Get(&id, query, user.Username, user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New("user not found")
	}

	return id, err
}

func (r *AuthDB) NewRefreshToken(token, userId string, expiresAt time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var existingTokenId string

	query := fmt.Sprintf(`
		SELECT ut.token_id
		FROM %s ut
		INNER JOIN %s t ON ut.token_id = t.id
		WHERE ut.user_id = $1`,
		usersTokensTable, usersTable,
	)
	err = tx.QueryRow(query, userId).Scan(&existingTokenId)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}

	if err == sql.ErrNoRows {
		// Insert Refresh Token
		query = fmt.Sprintf("INSERT INTO %s (refresh_token, expires_at) values ($1, $2) RETURNING id", tokensTable)
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
	} else {
		// Update Refresh Token
		query = fmt.Sprintf("UPDATE %s SET refresh_token = $1, expires_at = $2 WHERE id = $3", tokensTable)
		_, err = tx.Exec(query, token, expiresAt, existingTokenId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *AuthDB) CheckRefreshToken(userId, refreshToken string) error {
	var tokenId string
	query := fmt.Sprintf("SELECT ut.token_id FROM %s ut WHERE ut.user_id = $1", usersTokensTable)
	err := r.db.QueryRow(query, userId).Scan(&tokenId)
	if err != nil {
		return err
	}

	var storedToken string
	var expiresAt time.Time

	query = fmt.Sprintf("SELECT t.refresh_token, t.expires_at FROM %s t WHERE t.id = $1", tokensTable)
	err = r.db.QueryRow(query, tokenId).Scan(&storedToken, &expiresAt)
	if err != nil {
		return err
	}

	if storedToken != refreshToken {
		return errors.New("tokens are different")
	}

	if time.Now().After(expiresAt) {
		return errors.New("token has expired")
	}

	return nil
}

func (r *AuthDB) DeleteRefreshToken(userId string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s t
		USING %s ut
		WHERE t.id = ut.token_id AND ut.user_id = $1`,
		tokensTable, usersTokensTable,
	)
	_, err := r.db.Exec(query, userId)

	return err
}
