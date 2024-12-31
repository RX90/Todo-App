package service

import (
	"github.com/RX90/Todo-App/server/internal/repository"
	"github.com/RX90/Todo-App/server/internal/todo"
)

type Authorization interface {
	CreateUser(user todo.User) error
	GetUserId(user todo.User) (string, error)
	NewAccessToken(userId string) (string, error)
	NewRefreshToken(userId string) (string, error)
	ParseAccessToken(token string) (string, error)
	CheckRefreshToken(userId, refreshToken string) error
}

type TodoList interface {
	Create(userId string, list todo.List) error
}

type Service struct {
	Authorization
	TodoList
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: newAuthService(repos.Authorization),
		TodoList:      newTodoListService(repos.TodoList),
	}
}
