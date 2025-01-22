package service

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/RX90/Todo-App/server/internal/repository"
	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "f3by1efb08y1f0b8"
	signingKey = "vr3urn93u1dnwi00"
	accessTTL  = 30 * time.Minute    // 30 Minutes
	RefreshTTL = 30 * 24 * time.Hour // 30 Days
)

type AuthService struct {
	repos repository.Authorization
}

type TokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

func newAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) CreateUser(user todo.User) error {
	user.Password = generatePasswordHash(user.Password)
	return s.repos.CreateUser(user)
}

func (s *AuthService) GetUserId(user todo.User) (string, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repos.GetUserId(user)
}

func (s *AuthService) NewAccessToken(userId string) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTTL).Unix(),
			Subject:   userId,
		},
		userId,
	})

	tokenString, err := accessToken.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) NewRefreshToken(userId string) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := fmt.Sprintf("%x", b)
	expiresAt := time.Now().Add(RefreshTTL)

	if err := s.repos.NewRefreshToken(token, userId, expiresAt); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ParseAccessToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&TokenClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(signingKey), nil
		})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				if claims, ok := token.Claims.(*TokenClaims); ok {
					return claims.UserId, errors.New("token has expired")
				}
			}
		}
		return "", err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *TokenClaims")
	}
	if !token.Valid {
		return "", errors.New("token is invalid")
	}

	return claims.UserId, nil
}

func (s *AuthService) CheckRefreshToken(userId, refreshToken string) error {
	return s.repos.CheckRefreshToken(userId, refreshToken)
}

func (s *AuthService) DeleteRefreshToken(userId, refreshToken string) error {
	return s.repos.DeleteRefreshToken(userId, refreshToken)
}
