package handler

import (
	"net/http"
	"strconv"
	"time"

	todo "github.com/RX90/Todo-App"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, userId, err := h.services.Authorization.NewAccessToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "User does not exist")
		return
	}

	refreshToken, err := h.services.Authorization.NewRefreshToken()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	exp := time.Now().Add(15 * time.Minute)

	id, err := h.services.Authorization.CreateToken(refreshToken, exp, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  exp,
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, cookie)

	cookie = &http.Cookie{
		Name:     "userId",
		Value:    strconv.Itoa(userId),
		Path:     "/",
		Domain:   "localhost",
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{
		"refresh_id": id,
		"token":      accessToken,
	})
}

func (h *Handler) refresh(c *gin.Context) {
	// ЧЕ ДЕЛАТЬ??????????????? КАК ПРОВЕРЯТЬ REFRESH ТОКЕН?????????????
}
