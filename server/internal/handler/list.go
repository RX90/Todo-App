package handler

import (
	"net/http"

	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	var input todo.List

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := h.services.TodoList.Create(userId, input); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
