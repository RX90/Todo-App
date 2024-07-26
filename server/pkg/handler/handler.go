package handler

import (
	"net/http"

	"github.com/RX90/Todo-App/pkg/service"
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
	router.LoadHTMLGlob("../client/templates/*.html")

	router.GET("/", func(c *gin.Context) {})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)

		auth.GET("/sign-up", func(c *gin.Context) {
			c.HTML(http.StatusOK, "sign-up.html", nil)
		})
		auth.GET("/sign-in", func(c *gin.Context) {
			c.HTML(http.StatusOK, "sign-in.html", nil)
		})
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}
		items := api.Group("items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}

	return router
}
