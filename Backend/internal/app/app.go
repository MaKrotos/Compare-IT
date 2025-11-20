package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/router"
)

type App struct {
	server *http.Server
}

func New() *App {
	// Создаем маршрутизатор
	r := router.New()

	// Добавляем middleware
	r.Use(middleware.Logging)
	r.Use(middleware.ErrorHandler)

	// Регистрируем обработчики
	handlers.RegisterRoutes(r)

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	return &App{
		server: server,
	}
}

func (a *App) Run() {
	// Запуск сервера в отдельной горутине
	go func() {
		fmt.Println("Server is running on port 8080...")
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Server is shutting down...")

	// Создаем контекст с таймаутом для корректного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Пытаемся корректно завершить сервер
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}

	fmt.Println("Server stopped")
}