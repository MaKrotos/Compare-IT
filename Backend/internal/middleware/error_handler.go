package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

func ErrorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Создаем ResponseWriter, который отслеживает статус код
		hw := &hookedResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		// Вызываем следующий обработчик
		next(hw, r)
		
		// Если произошла ошибка (код >= 400), логируем её
		if hw.statusCode >= http.StatusBadRequest {
			log.Printf("HTTP %d %s %s", hw.statusCode, r.Method, r.URL.Path)
		}
	}
}

type hookedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (hrw *hookedResponseWriter) WriteHeader(code int) {
	hrw.statusCode = code
	hrw.ResponseWriter.WriteHeader(code)
}

// Error представляет стандартную структуру ошибки API
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SendError отправляет стандартный JSON ответ об ошибке
func SendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	errorResponse := Error{
		Message: message,
		Code:    code,
	}
	
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Failed to send error response: %v", err)
	}
}