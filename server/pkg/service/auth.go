package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	todo "github.com/RX90/Todo-App"
	"github.com/RX90/Todo-App/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "hjqrhjqw124617ajfhajs1231231ckmd"
	signingKey = "qrkjk#4jdaDDSFJlja#4353KSFjH!3ki"
	tokenTTL   = 5 * time.Minute
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) NewAccessToken(username, password string) (string, int, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", 0, err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			Subject:   strconv.Itoa(user.Id),
		},
		user.Id,
	})

	tokenString, err := accessToken.SignedString([]byte(signingKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, user.Id, nil
}

func (s *AuthService) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	src := rand.NewSource(time.Now().UnixMilli())
	r := rand.New(src)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(signingKey), nil
		},
	)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *AuthService) CreateToken(token string, exp time.Time, userId int) (int, error) {
	return s.repo.CreateToken(token, exp, userId)
}

func (s *AuthService) DeleteToken(userId string, token string) (error) {
	return s.repo.DeleteToken(userId, token)
}
