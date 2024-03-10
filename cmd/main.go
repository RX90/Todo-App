package main

import (
	"log"

	todo "github.com/RX90/Todo-App"
	"github.com/RX90/Todo-App/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)
	srv := new(todo.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
} 	