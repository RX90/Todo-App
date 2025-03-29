package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeader = "Authorization"
	userCtx    = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "auth header is empty"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" || headerParts[1] == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "auth header is invalid"})
		return
	}

	userId, err := h.services.Authorization.ParseAccessToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	c.Set(userCtx, userId)
	c.Next()
}

func getUserCtx(c *gin.Context) string {
	idAny, _ := c.Get(userCtx)

	return idAny.(string)
}

func authInputValidation(username, password string) error {
	l1, l2 := len(username), len(password)
	hasLetter, hasDigit := false, false

	switch {
	case l1 == 0:
		return errors.New("username is empty")
	case l2 == 0:
		return errors.New("password is empty")
	case l1 > 32:
		return errors.New("username exceeds 32 bytes")
	case l2 > 32:
		return errors.New("password exceeds 32 bytes")
	case l1 < 3:
		return errors.New("username is less than 3 characters")
	case l2 < 8:
		return errors.New("password is less than 8 characters")
	}

	for _, char := range username {
		if !isAllowed(char) {
			return errors.New("username has invalid character")
		}
	}

	for _, char := range password {
		if !isAllowed(char) {
			return errors.New("password has invalid character")
		}

		if isLetter(char) {
			hasLetter = true
		} else if isDigit(char) {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return errors.New("password must contain at least one english letter and one digit")
	}

	return nil
}

func isAllowed(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') || r == '-' || r == '_'
}

func isLetter(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}
