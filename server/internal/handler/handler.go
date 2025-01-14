package handler

import (
	"net/http"

	"github.com/RX90/Todo-App/server/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Static("/static", "../client/static")
	router.Static("/src", "../client/src")
	router.LoadHTMLGlob("../client/templates/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main.html", nil)
	})

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.signIn)
			auth.POST("/refresh", h.refreshTokens)
			auth.POST("/logout", h.userIdentity, h.logout)
		}

		lists := api.Group("/lists", h.userIdentity)
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.PUT("/:listId", h.updateList)
			lists.DELETE("/:listId", h.deleteList)

			tasks := lists.Group(":listId/tasks")
			{
				tasks.POST("/", h.createTask)
				tasks.GET("/", h.getAllTasks)
				tasks.PUT("/:taskId", h.updateTask)
				tasks.DELETE("/:taskId", h.deleteTask)
			}
		}
	}

	return router
}
