package service

import "github.com/RX90/Todo-App/server/internal/repository"

type AuthService struct {
	repos repository.Authorization
}

func newAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}
