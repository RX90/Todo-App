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
	tasksTable       = "tasks"
	listsTasksTable  = "lists_tasks"
	tokensTable      = "tokens"
	usersTokensTable = "users_tokens"
)

type Authorization interface {
	CreateUser(user todo.User) error
	GetUserId(user todo.User) (string, error)
	NewRefreshToken(token, userId string, expiresAt time.Time) error
	CheckRefreshToken(userId, refreshToken string) error
	DeleteRefreshToken(userId, refreshToken string) error
}

type TodoList interface {
	Create(userId string, list todo.List) (string, error)
	GetAll(userId string) ([]todo.List, error)
	GetById(userId, listId string) (todo.List, error)
	Update(userId, listId string, list todo.List) error
	Delete(userId, listId string) error
}

type TodoTask interface {
	Create(listId string, task todo.Task) (string, error)
	GetAll(userId, listId string) ([]todo.Task, error)
	Update(userId, listId, taskId string, task todo.UpdateTaskInput) error
	Delete(userId, listId, taskId string) error
}

type Repository struct {
	Authorization
	TodoList
	TodoTask
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthDB(db),
		TodoList:      newListDB(db),
		TodoTask:      newTaskDB(db),
	}
}
