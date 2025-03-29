package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createTask(c *gin.Context) {
	userId := getUserCtx(c)

	listId := c.Param("listId")
	_, err := strconv.Atoi(listId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't get list id: %s", err.Error())})
		return
	}

	var input todo.Task

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't bind JSON: %s", err.Error())})
		return
	}

	l := utf8.RuneCountInString(input.Title)

	if l == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "task title is empty"})
		return
	} else if l > 255 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "task title exceeds 255 characters"})
		return
	}

	taskId, err := h.services.TodoTask.Create(userId, listId, input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't create task: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task_id": taskId})
}

func (h *Handler) getAllTasks(c *gin.Context) {
	userId := getUserCtx(c)

	listId := c.Param("listId")
	_, err := strconv.Atoi(listId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't get list id: %s", err.Error())})
		return
	}

	tasks, err := h.services.TodoTask.GetAll(userId, listId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't get all tasks: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) updateTask(c *gin.Context) {
	userId := getUserCtx(c)

	listId := c.Param("listId")
	_, err := strconv.Atoi(listId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't get list id: %s", err.Error())})
		return
	}

	taskId := c.Param("taskId")
	_, err = strconv.Atoi(taskId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't get task id: %s", err.Error())})
		return
	}

	var input todo.UpdateTaskInput

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't bind JSON: %s", err.Error())})
		return
	}

	if err := input.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if input.Title != nil {
		l := utf8.RuneCountInString(*input.Title)

		if l == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "task title is empty"})
			return
		} else if l > 255 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "task title exceeds 255 characters"})
			return
		}
	}

	if err := h.services.TodoTask.Update(userId, listId, taskId, input); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't update task: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) deleteTask(c *gin.Context) {
	userId := getUserCtx(c)

	listId := c.Param("listId")
	_, err := strconv.Atoi(listId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't get list id: %s", err.Error())})
		return
	}

	taskId := c.Param("taskId")
	_, err = strconv.Atoi(taskId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't get task id: %s", err.Error())})
		return
	}

	err = h.services.TodoTask.Delete(userId, listId, taskId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't delete task: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
