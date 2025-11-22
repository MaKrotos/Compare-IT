package handlers

import (
	"encoding/json"
	"net/http"

	"backend/internal/database"
	"backend/internal/middleware"
)

// GeneratePublicLinkHandler создает публичную ссылку для коллекции
func GeneratePublicLinkHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем информацию о пользователе из контекста
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		middleware.SendError(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// Декодируем тело запроса
	var requestData struct {
		CollectionID int `json:"collection_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		middleware.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем, что ID коллекции передан
	if requestData.CollectionID == 0 {
		middleware.SendError(w, "Collection ID is required", http.StatusBadRequest)
		return
	}

	// Сначала пытаемся получить коллекцию сравнений
	comparisonCollection, err := database.GetComparisonCollectionByID(requestData.CollectionID, userID)
	if err != nil {
		// Если произошла ошибка при получении коллекции сравнений, пробуем получить обычную коллекцию
		collection, err := database.GetCollectionByID(requestData.CollectionID, userID)
		if err != nil {
			middleware.SendError(w, "Failed to get collection: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		if collection == nil {
			middleware.SendError(w, "Collection not found", http.StatusNotFound)
			return
		}
		
		// Генерируем публичную ссылку для обычной коллекции
		publicLink, err := database.GeneratePublicLinkForCollection(requestData.CollectionID)
		if err != nil {
			middleware.SendError(w, "Failed to generate public link: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Отправляем ответ
		response := map[string]interface{}{
			"public_link": publicLink,
			"message":    "Public link generated successfully",
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Если получили коллекцию сравнений
	if comparisonCollection != nil {
		// Генерируем публичную ссылку для коллекции сравнений
		publicLink, err := database.GeneratePublicLink(requestData.CollectionID)
		if err != nil {
			middleware.SendError(w, "Failed to generate public link: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Отправляем ответ
		response := map[string]interface{}{
			"public_link": publicLink,
			"message":    "Public link generated successfully",
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Если не нашли ни одну коллекцию, пытаемся получить обычную коллекцию
	collection, err := database.GetCollectionByID(requestData.CollectionID, userID)
	if err != nil {
		middleware.SendError(w, "Failed to get collection: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if collection == nil {
		middleware.SendError(w, "Collection not found", http.StatusNotFound)
		return
	}
	
	// Генерируем публичную ссылку для обычной коллекции
	publicLink, err := database.GeneratePublicLinkForCollection(requestData.CollectionID)
	if err != nil {
		middleware.SendError(w, "Failed to generate public link: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Отправляем ответ
	response := map[string]interface{}{
		"public_link": publicLink,
		"message":    "Public link generated successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RemovePublicLinkHandler удаляет публичную ссылку для коллекции
func RemovePublicLinkHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем информацию о пользователе из контекста
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		middleware.SendError(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// Декодируем тело запроса
	var requestData struct {
		CollectionID int `json:"collection_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		middleware.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем, что ID коллекции передан
	if requestData.CollectionID == 0 {
		middleware.SendError(w, "Collection ID is required", http.StatusBadRequest)
		return
	}

	// Сначала пытаемся получить коллекцию сравнений
	comparisonCollection, err := database.GetComparisonCollectionByID(requestData.CollectionID, userID)
	if err != nil {
		// Если произошла ошибка при получении коллекции сравнений, пробуем получить обычную коллекцию
		collection, err := database.GetCollectionByID(requestData.CollectionID, userID)
		if err != nil {
			middleware.SendError(w, "Failed to get collection: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		if collection == nil {
			middleware.SendError(w, "Collection not found", http.StatusNotFound)
			return
		}
		
		// Удаляем публичную ссылку для обычной коллекции
		err = database.RemovePublicLinkForCollection(requestData.CollectionID)
		if err != nil {
			middleware.SendError(w, "Failed to remove public link: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Отправляем ответ
		response := map[string]interface{}{
			"message": "Public link removed successfully",
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Если получили коллекцию сравнений
	if comparisonCollection != nil {
		// Удаляем публичную ссылку для коллекции сравнений
		err = database.RemovePublicLink(requestData.CollectionID)
		if err != nil {
			middleware.SendError(w, "Failed to remove public link: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Отправляем ответ
		response := map[string]interface{}{
			"message": "Public link removed successfully",
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Если не нашли ни одну коллекцию, пытаемся получить обычную коллекцию
	collection, err := database.GetCollectionByID(requestData.CollectionID, userID)
	if err != nil {
		middleware.SendError(w, "Failed to get collection: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if collection == nil {
		middleware.SendError(w, "Collection not found", http.StatusNotFound)
		return
	}
	
	// Удаляем публичную ссылку для обычной коллекции
	err = database.RemovePublicLinkForCollection(requestData.CollectionID)
	if err != nil {
		middleware.SendError(w, "Failed to remove public link: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Отправляем ответ
	response := map[string]interface{}{
		"message": "Public link removed successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetPublicCollectionHandler получает коллекцию по публичной ссылке
func GetPublicCollectionHandler(w http.ResponseWriter, r *http.Request) {
	// Декодируем тело запроса
	var requestData struct {
		PublicLink string `json:"public_link"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		middleware.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем, что публичная ссылка передана
	if requestData.PublicLink == "" {
		middleware.SendError(w, "Public link is required", http.StatusBadRequest)
		return
	}

	// Сначала пытаемся получить коллекцию сравнений по публичной ссылке
	comparisonCollection, err := database.GetComparisonCollectionByPublicLink(requestData.PublicLink)
	if err != nil {
		middleware.SendError(w, "Failed to get collection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Если нашли коллекцию сравнений, возвращаем её
	if comparisonCollection != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comparisonCollection)
		return
	}

	// Если не нашли коллекцию сравнений, пытаемся получить обычную коллекцию
	collection, err := database.GetCollectionByPublicLink(requestData.PublicLink)
	if err != nil {
		middleware.SendError(w, "Failed to get collection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Если не нашли ни одну коллекцию, возвращаем ошибку
	if collection == nil {
		middleware.SendError(w, "Collection not found", http.StatusNotFound)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(collection)
}