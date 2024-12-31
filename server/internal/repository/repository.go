package repository

import (
	"time"

	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable       = "users"
	listsTable       = "lists"
	usersListsTable  = "users_lists"
	itemsTable       = "items"
	listsItemsTable  = "lists_items"
	tokensTable      = "tokens"
	usersTokensTable = "users_tokens"
)

type Authorization interface {
	CreateUser(user todo.User) error
	GetUserId(user todo.User) (string, error)
	NewRefreshToken(token, userId string, expiresAt time.Time) error
	CheckRefreshToken(userId, refreshToken string) error
}

type TodoList interface {
	Create(userId string, list todo.List) error
}

type Repository struct {
	Authorization
	TodoList
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthDB(db),
		TodoList:      newListDB(db),
	}
}
