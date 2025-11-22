package config

import (
	"os"
	"sync"
)

type Config struct {
	JWTSecret         string
	TelegramBotToken  string
	TelegramWebhookURL string
}

var (
	config *Config
	once   sync.Once
)

// GetConfig возвращает синглтон конфигурации
func GetConfig() *Config {
	once.Do(func() {
		config = &Config{
			JWTSecret:          getEnv("JWT_SECRET", "default_secret_key"),
			TelegramBotToken:  getEnv("TELEGRAM_BOT_TOKEN", ""),
			TelegramWebhookURL: getEnv("TELEGRAM_WEBHOOK_URL", ""),
		}
	})
	return config
}

// getEnv получает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}