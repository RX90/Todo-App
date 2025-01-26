package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/RX90/Todo-App/server/internal/service"
	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/gin-gonic/gin"
)

const (
	refresh = "refreshToken"
)

type Response struct {
	Message string `json:"message"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't bind JSON: %s", err.Error())})
		return
	}

	if err := inputValidate(input.Username, input.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if utf8.RuneCountInString(input.Password) < 8 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "password is less than 8 characters"})
		return
	}

	if err := h.services.Authorization.CreateUser(input); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't create user: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) signIn(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't bind JSON: %s", err.Error())})
		return
	}

	if err := inputValidate(input.Username, input.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	userId, err := h.services.Authorization.GetUserId(input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": fmt.Sprintf("can't get user id: %s", err.Error())})
		return
	}

	accessToken, err := h.services.Authorization.NewAccessToken(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't create access token: %s", err.Error())})
		return
	}

	refreshToken, err := h.services.Authorization.NewRefreshToken(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't create refresh token: %s", err.Error())})
		return
	}

	cookie := &http.Cookie{
		Name:     refresh,
		Value:    refreshToken,
		Expires:  time.Now().Add(service.RefreshTTL),
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (h *Handler) refreshTokens(c *gin.Context) {
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

	accessToken := headerParts[1]
	refreshToken, err := c.Cookie(refresh)
	if err != nil || refreshToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "refresh token is missing"})
		return
	}

	userId, err := h.services.Authorization.ParseAccessToken(accessToken)
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": fmt.Sprintf("access token is invalid: %s", err.Error())})
		return
	}

	if err := h.services.Authorization.CheckRefreshToken(userId, refreshToken); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": fmt.Sprintf("refresh token is invalid: %s", err.Error())})
		return
	}

	accessToken, err = h.services.Authorization.NewAccessToken(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't create access token: %s", err.Error())})
		return
	}

	refreshToken, err = h.services.Authorization.NewRefreshToken(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't create refresh token: %s", err.Error())})
		return
	}

	cookie := &http.Cookie{
		Name:     refresh,
		Value:    refreshToken,
		Expires:  time.Now().Add(service.RefreshTTL),
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (h *Handler) logout(c *gin.Context) {
	userId := getUserCtx(c)

	refreshToken, err := c.Cookie(refresh)
	if err != nil || refreshToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "refresh token is missing"})
		return
	}

	if err := h.services.Authorization.DeleteRefreshToken(userId, refreshToken); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't delete refresh token: %s", err.Error())})
		return
	}

	cookie := &http.Cookie{
		Name:   refresh,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
