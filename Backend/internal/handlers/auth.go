package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"backend/auth"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"
)


// TelegramAuthRequest представляет запрос на аутентификацию через Telegram
type TelegramAuthRequest struct {
	InitData string `json:"initData"`
}

// TelegramAuthResponse представляет ответ на аутентификацию через Telegram
type TelegramAuthResponse struct {
	Token string `json:"token"`
	User  *models.TelegramUser `json:"user"`
}

// getUserFromContext извлекает информацию о пользователе из контекста запроса
func getUserFromContext(ctx context.Context) (*models.TelegramUser, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	
	telegramID, ok := ctx.Value("telegram_id").(int64)
	if !ok {
		return nil, fmt.Errorf("telegram ID not found in context")
	}
	
	generatedName, ok := ctx.Value("generated_name").(string)
	if !ok {
		return nil, fmt.Errorf("generated name not found in context")
	}
	
	isAdmin, ok := ctx.Value("is_admin").(bool)
	if !ok {
		return nil, fmt.Errorf("admin status not found in context")
	}
	
	role, ok := ctx.Value("role").(models.UserRole)
	if !ok {
		return nil, fmt.Errorf("role not found in context")
	}

	// Создаем объект пользователя
	user := &models.TelegramUser{
		ID:            userID,
		TelegramID:    telegramID,
		GeneratedName: generatedName,
		IsAdmin:       isAdmin,
		Role:          role,
	}

	return user, nil
}

// authenticateTelegramUser выполняет аутентификацию пользователя через Telegram
func authenticateTelegramUser(initData string) (*TelegramAuthResponse, error) {
	// Получаем конфигурацию
	cfg := config.GetConfig()
	
	// Создаем валидатор Telegram
	telegramAuth := auth.NewTelegramAuth(cfg.TelegramBotToken)

	// Валидируем initData
	isValid, telegramData, err := telegramAuth.ValidateInitData(initData)
	if err != nil {
		return nil, fmt.Errorf("invalid Telegram initData: %w", err)
	}
	
	if !isValid {
		return nil, fmt.Errorf("invalid Telegram initData")
	}

	// Проверяем, существует ли пользователь в базе данных
	user, err := database.GetUserByTelegramID(telegramData.ID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Если пользователь не существует, создаем его
	if user == nil {
		user = &models.TelegramUser{
			TelegramID:    telegramData.ID,
			FirstName:     telegramData.FirstName,
			LastName:      telegramData.LastName,
			Username:      telegramData.Username,
			GeneratedName: telegramData.FirstName, // В реальном приложении здесь должна быть генерация случайного имени
			IsAdmin:       false, // Установите в true для администраторов
			Role:          models.RoleUser,
		}
		
		// Сохраняем пользователя в базе данных
		if err := database.CreateUser(user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Обновляем время последней активности пользователя
	lastActive := time.Now().Format(time.RFC3339)
	if err := database.UpdateUserLastActive(user.ID, lastActive); err != nil {
		// Логируем ошибку, но не прерываем процесс аутентификации
		// В реальном приложении здесь должна быть настоящая логика логирования
	}

	// Генерируем JWT токен
	token, err := auth.GenerateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Подготавливаем ответ
	response := &TelegramAuthResponse{
		Token: token,
		User:  user,
	}

	return response, nil
}

// TelegramAuthHandler обрабатывает аутентификацию через Telegram WebApp
func TelegramAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS для всех origins
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Обрабатываем preflight запросы
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Проверяем метод запроса
	if r.Method != "POST" {
		middleware.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем тело запроса
	var req TelegramAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Выполняем аутентификацию через Telegram
	response, err := authenticateTelegramUser(req.InitData)
	if err != nil {
		middleware.SendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ProfileHandler возвращает информацию о текущем пользователе
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем информацию о пользователе из контекста
	user, err := getUserFromContext(r.Context())
	if err != nil {
		middleware.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}