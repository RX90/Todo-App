package repository

import (
	"fmt"
	"time"

	todo "github.com/RX90/Todo-App"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	var user todo.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}

func (r *AuthPostgres) CreateToken(token string, exp time.Time, userId int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var tokenId int

	query := fmt.Sprintf("INSERT INTO %s (refresh_token, expires_in) values ($1, $2) RETURNING id", tokensTable)
	row := tx.QueryRow(query, token, exp)
	if err := row.Scan(&tokenId); err != nil {
		tx.Rollback()
		return 0, err
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, token_id) values ($1, $2)", usersTokensTable)
	_, err = tx.Exec(query, userId, tokenId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return tokenId, tx.Commit()
}
