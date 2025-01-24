package handler

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

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

func inputValidate(input ...string) error {
	for _, data := range input {
		l := utf8.RuneCountInString(data)
		
		if l == 0 {
			return errors.New("empty input")
		} else if l > 32 {
			return errors.New("input exceeds 32 characters")
		}
	}

	return nil
}