package repository

import (
	"time"

	"github.com/RX90/Todo-App/server/internal/user"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable       = "users"
	todoListsTable   = "todo_lists"
	usersListsTable  = "users_lists"
	todoItemsTable   = "todo_items"
	listsItemsTable  = "lists_items"
	tokensTable      = "tokens"
	usersTokensTable = "users_tokens"
)

type Authorization interface {
	CreateUser(user user.User) error
	GetUserId(username, password string) (string, error)
	NewRefreshToken(token, userId string, expiresAt time.Time) error
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthDB(db),
	}
}
