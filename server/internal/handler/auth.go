package handler

import (
	"net/http"

	"github.com/RX90/Todo-App/server/internal/user"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input user.User

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := h.services.Authorization.CreateUser(input); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
