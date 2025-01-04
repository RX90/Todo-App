package handler

import (
	"net/http"
	"strconv"

	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createTask(c *gin.Context) {
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

	var input todo.Task

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	taskId, err := h.services.TodoTask.Create(userId, listId, input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task_id": taskId})
}

func (h *Handler) getAllTasks(c *gin.Context) {
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

	tasks, err := h.services.TodoTask.GetAll(userId, listId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) updateTask(c *gin.Context) {
	userId, err := getUserCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	taskId := c.Param("id")
	_, err = strconv.Atoi(taskId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid id param"})
		return
	}

	var input todo.Task

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := h.services.TodoTask.Update(userId, taskId, input); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) deleteTask(c *gin.Context) {
	userId, err := getUserCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	taskId := c.Param("id")
	_, err = strconv.Atoi(taskId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid id param"})
		return
	}

	err = h.services.TodoTask.Delete(userId, taskId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
