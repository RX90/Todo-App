package service

import (
	"github.com/RX90/Todo-App/server/internal/repository"
	"github.com/RX90/Todo-App/server/internal/todo"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user todo.User) error
	GetUserId(user todo.User) (string, error)
	NewAccessToken(userId string) (string, error)
	NewRefreshToken(userId string) (string, error)
	ParseAccessToken(token string) (string, error)
	CheckRefreshToken(userId, refreshToken string) error
	DeleteRefreshToken(userId string) error
}

type TodoList interface {
	Create(userId string, list todo.List) (string, error)
	GetAll(userId string) ([]todo.List, error)
	Update(userId, listId string, list todo.List) error
	Delete(userId, listId string) error
}

type TodoTask interface {
	Create(userId, listId string, task todo.Task) (string, error)
	GetAll(userId, listId string) ([]todo.Task, error)
	Update(userId, listId, taskId string, task todo.UpdateTaskInput) error
	Delete(userId, listId, taskId string) error
}

type Service struct {
	Authorization
	TodoList
	TodoTask
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: newAuthService(repos.Authorization),
		TodoList:      newListService(repos.TodoList),
		TodoTask:      newTaskService(repos.TodoTask, repos.TodoList),
	}
}
