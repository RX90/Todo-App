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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "access token is empty"})
		return
	}
	
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "access token is invalid"})
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

func getUserId(c *gin.Context) (string, error) {
	idAny, ok := c.Get(userCtx)
	if !ok {
		return "", errors.New("user id not found")
	}

	id, ok := idAny.(string)
	if !ok {
		return "", errors.New("user id not found")
	}

	return id, nil
}
