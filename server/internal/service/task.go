package service

import (
	"github.com/RX90/Todo-App/server/internal/repository"
	"github.com/RX90/Todo-App/server/internal/todo"
)

type TaskService struct {
	repos     repository.TodoTask
	listRepos repository.TodoList
}

func newTaskService(repos repository.TodoTask, listRepos repository.TodoList) *TaskService {
	return &TaskService{repos: repos, listRepos: listRepos}
}

func (s *TaskService) Create(userId, listId string, task todo.Task) (string, error) {
	_, err := s.listRepos.GetById(userId, listId)
	if err != nil {
		return "", err
	}

	return s.repos.Create(listId, task)
}

func (s *TaskService) GetAll(userId, listId string) ([]todo.Task, error) {
	return s.repos.GetAll(userId, listId)
}

func (s *TaskService) Update(userId, listId, taskId string, task todo.UpdateTaskInput) error {
	return s.repos.Update(userId, listId, taskId, task)
}

func (s *TaskService) Delete(userId, listId, taskId string) error {
	return s.repos.Delete(userId, listId, taskId)
}
