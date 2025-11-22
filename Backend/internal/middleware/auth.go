package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/auth"
	"backend/internal/models"
)

// JWTAuthMiddleware middleware для проверки JWT токена
func JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			SendError(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Проверка формата заголовка
		if !strings.HasPrefix(authHeader, "Bearer ") {
			SendError(w, "Authorization header must start with Bearer", http.StatusUnauthorized)
			return
		}

		// Извлечение токена
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Валидация токена
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			SendError(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Сохранение claims в контексте запроса
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "telegram_id", claims.TelegramID)
		ctx = context.WithValue(ctx, "generated_name", claims.GeneratedName)
		ctx = context.WithValue(ctx, "is_admin", claims.IsAdmin)
		ctx = context.WithValue(ctx, "role", models.ParseRoleFromInt(claims.Role))
		ctx = context.WithValue(ctx, "claims", claims)
		
		// Передаем запрос с обновленным контекстом следующему обработчику
		next(w, r.WithContext(ctx))
	}
}

// AdminAuthMiddleware middleware для проверки прав администратора
func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение информации о правах администратора из контекста
		ctx := r.Context()
		isAdmin, ok := ctx.Value("is_admin").(bool)
		if !ok {
			SendError(w, "Admin privileges not found in context", http.StatusUnauthorized)
			return
		}

		// Проверка, является ли пользователь администратором
		if !isAdmin {
			SendError(w, "Admin privileges required", http.StatusForbidden)
			return
		}

		// Продолжение выполнения
		next(w, r)
	}
}