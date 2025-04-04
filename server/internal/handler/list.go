package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	userId := getUserCtx(c)

	var input todo.List

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't bind JSON: %s", err.Error())})
		return
	}

	l := utf8.RuneCountInString(input.Title)

	if l == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "list title is empty"})
		return
	} else if l > 32 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "list title exceeds 32 characters"})
		return
	}

	listId, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't create list: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list_id": listId})
}

func (h *Handler) getAllLists(c *gin.Context) {
	userId := getUserCtx(c)

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't get all lists: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, lists)
}

func (h *Handler) updateList(c *gin.Context) {
	userId := getUserCtx(c)

	listId := c.Param("listId")
	_, err := strconv.Atoi(listId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't get list id: %s", err.Error())})
		return
	}

	var input todo.List

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't bind JSON: %s", err.Error())})
		return
	}

	l := utf8.RuneCountInString(input.Title)

	if l == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "list title is empty"})
		return
	} else if l > 32 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "list title exceeds 32 characters"})
		return
	}

	if err := h.services.TodoList.Update(userId, listId, input); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't update list: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) deleteList(c *gin.Context) {
	userId := getUserCtx(c)

	listId := c.Param("listId")
	_, err := strconv.Atoi(listId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't get list id: %s", err.Error())})
		return
	}

	if err := h.services.TodoList.Delete(userId, listId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't delete list: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
