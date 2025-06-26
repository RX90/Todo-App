package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RX90/Todo-App/server"
	"github.com/RX90/Todo-App/server/internal/db"
	"github.com/RX90/Todo-App/server/internal/handler"
	"github.com/RX90/Todo-App/server/internal/repository"
	"github.com/RX90/Todo-App/server/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type App struct {
	server   *server.Server
	db       *sqlx.DB
	repos    *repository.Repository
	services *service.Service
	handlers *handler.Handler
}

func NewApp() (*App, error) {
	server := server.NewServer()

	if err := initConfig(); err != nil {
		return nil, err
	}

	if err := checkConfig(); err != nil {
		return nil, err
	}

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := checkEnv(); err != nil {
		return nil, err
	}

	db, err := db.NewSQLiteDB()
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	return &App{
		server:   server,
		db:       db,
		repos:    repos,
		services: services,
		handlers: handlers,
	}, nil
}

func (a *App) Run() error {
	go func() {
		a.server.Run(
			viper.GetString("server.port"),
			a.handlers.InitRoutes(),
		)
	}()

	log.Println("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("TodoApp Shutting Down")

	if err := a.server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	if err := a.db.Close(); err != nil {
		return fmt.Errorf("error occurred on db connection close: %s", err.Error())
	}

	return nil
}

func initConfig() error {
	viper.AddConfigPath("server/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func checkConfig() error {
    required := []string{
        "server.port",
    }
    missing := []string{}

    for _, key := range required {
        if viper.GetString(key) == "" {
            missing = append(missing, key)
        }
    }

    if len(missing) > 0 {
        return fmt.Errorf("missing required config values: %v", missing)
    }

    return nil
}

func checkEnv() error {
    required := []string{
        "SERVICE_SALT",
        "SERVICE_KEY",
    }
    missing := []string{}

    for _, key := range required {
        if os.Getenv(key) == "" {
            missing = append(missing, key)
        }
    }

    if len(missing) > 0 {
        return fmt.Errorf("missing required environment variables: %v", missing)
    }

    return nil
}