package service

import (
	"github.com/RX90/Todo-App/server/internal/repository"
	"github.com/RX90/Todo-App/server/internal/todo"
)

type ListService struct {
	repos repository.TodoList
}

func newListService(repos repository.TodoList) *ListService {
	return &ListService{repos: repos}
}

func (s *ListService) Create(userId string, list todo.List) (string, error) {
	return s.repos.Create(userId, list)
}

func (s *ListService) GetAll(userId string) ([]todo.List, error) {
	return s.repos.GetAll(userId)
}

func (s *ListService) Update(userId, listId string, list todo.List) error {
	return s.repos.Update(userId, listId, list)
}

func (s *ListService) Delete(userId, listId string) error {
	return s.repos.Delete(userId, listId)
}