package middleware

import (
	"context"
	"net/http"

	"backend/internal/database"
	"backend/internal/models"
)

// LoadUserComparisonsMiddleware загружает коллекции сравнений пользователя в контекст
func LoadUserComparisonsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем информацию о пользователе из контекста
		ctx := r.Context()
		userID, ok := ctx.Value("user_id").(int)
		if !ok {
			SendError(w, "User not found in context", http.StatusInternalServerError)
			return
		}

		// Получаем коллекции сравнений пользователя
		collections, err := database.GetComparisonCollectionsByUserID(userID)
		if err != nil {
			// Не прерываем выполнение, если не удалось получить коллекции
			// Просто логируем ошибку и продолжаем выполнение
			// В реальном приложении здесь должна быть настоящая логика логирования
			collections = []models.ComparisonCollection{}
		}

		// Добавляем коллекции в контекст
		ctx = context.WithValue(ctx, "user_comparisons", collections)
		
		// Передаем запрос с обновленным контекстом следующему обработчику
		next(w, r.WithContext(ctx))
	}
}