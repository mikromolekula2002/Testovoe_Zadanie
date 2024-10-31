package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/mikromolekula2002/Testovoe/internal/config"
	"github.com/mikromolekula2002/Testovoe/internal/controllers"
	"github.com/mikromolekula2002/Testovoe/internal/logger"
	"github.com/mikromolekula2002/Testovoe/internal/repo"
	"github.com/mikromolekula2002/Testovoe/internal/router"
	"github.com/mikromolekula2002/Testovoe/internal/services"
)

func main() {
	//Загрузка конфига
	cfg := config.LoadConfig("./config/config.yaml")
	// Инициализация базы данных
	logger := logger.Init(cfg.Logger.Level, cfg.Logger.FilePath, cfg.Logger.Output)

	//инициализация репозитория
	repo, err := repo.InitDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации Базы Данных: \n%v", err)
	}

	//инициализация сервисного слоя
	services := services.NewService(repo, cfg)
	newController := controllers.NewController(services, logger)

	// инициализация сервера
	router := router.NewRouter(cfg, newController)
	router.SetRoutes()

	newServer := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router.Gin,
	}

	// Start the server in a separate goroutine
	go func() {
		if err := newServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server start failed: %s\n", err)
		}
	}()
	log.Println("Server started successfully")

	// Set up signal catching
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown initiated...")

	// Context for graceful shutdown with a timeout of 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := newServer.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown Error: %v", err)
	}

	// Waiting for the shutdown context to be done or timeout
	select {
	case <-ctx.Done():
		log.Println("Server shutdown completed or timed out")
	}

	log.Println("Server exiting")
	os.Exit(0)
}
