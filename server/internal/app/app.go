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

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	db, err := db.NewPostgresDB(db.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
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
		// a.server.RunTLS(
		// 	viper.GetString("server.port"),
		// 	"server/certs/cert.crt",
		// 	"server/certs/cert.key",
		// 	a.handlers.InitRoutes(),
		// )
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
