package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/services/preview"
)

// PreviewHandler обрабатывает запросы на получение предварительного просмотра URL
func PreviewHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS для всех origins
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Обрабатываем preflight запросы
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	url := r.URL.Query().Get("url")
	if url == "" {
		middleware.SendError(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	previewData, err := preview.GetPreviewData(url)
	if err != nil {
		log.Printf("Error fetching preview data for URL %s: %v", url, err)
		// Вместо ошибки возвращаем минимальные данные
		previewData = &models.PreviewData{
			Title:       "Ошибка загрузки",
			Description: fmt.Sprintf("Не удалось загрузить данные: %v", err),
			Image:       "",
			URL:         url,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(previewData)
}