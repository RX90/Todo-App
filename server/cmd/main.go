package main

import (
	"log"

	"github.com/RX90/Todo-App/server/internal/app"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %s", err.Error())
	}

	if err := app.Run(); err != nil {
		log.Fatalf("App fatal error: %s", err.Error())
	}
}
