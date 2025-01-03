package handler

import (
	"net/http"
	"strconv"

	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	var input todo.List

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	listId, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list_id": listId})
}

func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lists)
}

func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	listId := c.Param("id")
	_, err = strconv.Atoi(listId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid id param"})
		return
	}

	var input todo.List

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := h.services.TodoList.Update(userId, listId, input); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	listId := c.Param("id")
	_, err = strconv.Atoi(listId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid id param"})
		return
	}

	if err := h.services.TodoList.Delete(userId, listId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
