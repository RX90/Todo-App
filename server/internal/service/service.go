package service

import "github.com/RX90/Todo-App/server/internal/repository"

type Authorization interface{}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: newAuthService(repos.Authorization),
	}
}
