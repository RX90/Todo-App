package service

import (
	"crypto/sha256"
	"fmt"

	"github.com/RX90/Todo-App/server/internal/repository"
	"github.com/RX90/Todo-App/server/internal/user"
)

const (
	salt = "f3by1efb08y1f0b8"
)

type AuthService struct {
	repos repository.Authorization
}

func newAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func isPasswordOK(password string) error {
	if len(password) < 8 || len(password) > 32 {
		return fmt.Errorf("invalid length of password")
	}
	var (
		hasDigit, hasLowerCase, hasUpperCase bool
	)

	for _, char := range password {
		if 'a' <= char && char <= 'z' {
			hasLowerCase = true
		} else if 'A' <= char && char <= 'Z' {
			hasUpperCase = true
		} else if '0' <= char && char <= '9' {
			hasDigit = true
		} else {
			return fmt.Errorf("invalid character used in password: '%s'", string(char))
		}
	}

	if !hasDigit {
		return fmt.Errorf("no digits in password")
	} else if !hasLowerCase {
		return fmt.Errorf("no lowercase letters in password")
	} else if !hasUpperCase {
		return fmt.Errorf("no uppercase letters in password")
	}

	return nil
}

func (s *AuthService) CreateUser(user user.User) error {
	if err := isPasswordOK(user.Password); err != nil {
		return err
	}

	user.Password = generatePasswordHash(user.Password)
	return s.repos.CreateUser(user)
}
