package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// ErrNoRows возвращает ошибку, когда запись не найдена
var ErrNoRows = errors.New("no rows found")

// InitDB инициализирует подключение к базе данных
func InitDB() {
	// Получаем строку подключения из переменной окружения
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Открываем подключение к базе данных
	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to database")
}

// GetDB возвращает экземпляр базы данных
func GetDB() *sql.DB {
	return db
}

// CloseDB закрывает подключение к базе данных
func CloseDB() {
	if db != nil {
		db.Close()
	}
}