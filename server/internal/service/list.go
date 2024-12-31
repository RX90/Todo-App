package service

import (
	"github.com/RX90/Todo-App/server/internal/repository"
	"github.com/RX90/Todo-App/server/internal/todo"
)

type TodoListService struct {
	repos repository.TodoList
}

func newTodoListService(repos repository.TodoList) *TodoListService {
	return &TodoListService{repos: repos}
}

func (s *TodoListService) Create(userId string, list todo.List) error {
	return s.repos.Create(userId, list)
}