package handlers

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "strconv"

  "backend/internal/database"
  "backend/internal/middleware"
  "backend/internal/models"
)

// CreateComparisonCollectionHandler создает новую коллекцию сравнений
func CreateComparisonCollectionHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем информацию о пользователе из контекста
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		middleware.SendError(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// Декодируем тело запроса
	var collection models.ComparisonCollection
	if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
		middleware.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Устанавливаем ID пользователя
	collection.UserID = userID
	
	// Устанавливаем значения по умолчанию для весов рейтингов, если они не заданы
	if collection.PriceRatingWeight == 0 && collection.ProsConsRatingWeight == 0 {
		collection.PriceRatingWeight = 20
		collection.ProsConsRatingWeight = 80
	}

	// Создаем коллекцию в базе данных
	if err := database.CreateComparisonCollection(&collection); err != nil {
		middleware.SendError(w, "Failed to create comparison collection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(collection)
}

// GetComparisonCollectionsHandler получает все коллекции сравнений пользователя
func GetComparisonCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем информацию о пользователе из контекста
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		middleware.SendError(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// Получаем коллекции из базы данных
	collections, err := database.GetComparisonCollectionsByUserID(userID)
	if err != nil {
		middleware.SendError(w, "Failed to get comparison collections: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(collections)
}

// UpdateComparisonCollectionHandler обновляет коллекцию сравнений
func UpdateComparisonCollectionHandler(w http.ResponseWriter, r *http.Request) {
  // Получаем информацию о пользователе из контекста
  ctx := r.Context()
  userID, ok := ctx.Value("user_id").(int)
  if !ok {
    middleware.SendError(w, "User not found in context", http.StatusInternalServerError)
    return
  }

  // Декодируем тело запроса
  var requestData struct {
    ID                   int                      `json:"id"`
    Name                 string                   `json:"name"`
    Items                []models.ComparisonItem  `json:"items"`
    PriceRatingWeight    int                      `json:"price_rating_weight"`
    ProsConsRatingWeight int                      `json:"pros_cons_rating_weight"`
  }
  if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
    middleware.SendError(w, "Invalid request body", http.StatusBadRequest)
    return
  }

  // Проверяем, что ID передан
  if requestData.ID == 0 {
    middleware.SendError(w, "Collection ID is required", http.StatusBadRequest)
    return
  }

  // Создаем объект коллекции
  collection := models.ComparisonCollection{
    ID:                   requestData.ID,
    UserID:               userID,
    Name:                 requestData.Name,
    Items:                requestData.Items,
    PriceRatingWeight:    requestData.PriceRatingWeight,
    ProsConsRatingWeight: requestData.ProsConsRatingWeight,
  }

  // Обновляем коллекцию в базе данных
  if err := database.UpdateComparisonCollection(&collection); err != nil {
    if err == database.ErrNoRows {
      middleware.SendError(w, "Comparison collection not found", http.StatusNotFound)
      return
    }
    middleware.SendError(w, "Failed to update comparison collection: "+err.Error(), http.StatusInternalServerError)
    return
  }

  // Отправляем ответ
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(collection)
}

// DeleteComparisonCollectionHandler удаляет коллекцию сравнений
func DeleteComparisonCollectionHandler(w http.ResponseWriter, r *http.Request) {
  // Получаем информацию о пользователе из контекста
  ctx := r.Context()
  userID, ok := ctx.Value("user_id").(int)
  if !ok {
    middleware.SendError(w, "User not found in context", http.StatusInternalServerError)
    return
  }

  // Декодируем тело запроса
  var requestData struct {
    ID int `json:"id"`
  }
  if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
    middleware.SendError(w, "Invalid request body", http.StatusBadRequest)
    return
  }

  // Проверяем, что ID передан
  if requestData.ID == 0 {
    middleware.SendError(w, "Collection ID is required", http.StatusBadRequest)
    return
  }

  // Удаляем коллекцию из базы данных
  if err := database.DeleteComparisonCollection(requestData.ID, userID); err != nil {
    if err == database.ErrNoRows {
      middleware.SendError(w, "Comparison collection not found", http.StatusNotFound)
      return
    }
    middleware.SendError(w, "Failed to delete comparison collection: "+err.Error(), http.StatusInternalServerError)
    return
  }

  // Отправляем ответ
  w.WriteHeader(http.StatusNoContent)
}

// GetComparisonCollectionByIdHandler получает коллекцию сравнений по ID
func GetComparisonCollectionByIdHandler(w http.ResponseWriter, r *http.Request) {
  // Получаем информацию о пользователе из контекста
  ctx := r.Context()
  userID, ok := ctx.Value("user_id").(int)
  if !ok {
    middleware.SendError(w, "User not found in context", http.StatusInternalServerError)
    return
  }

  // Декодируем тело запроса
  var requestData struct {
    ID interface{} `json:"id"`
  }
  
  // Логируем тело запроса для отладки
  bodyBytes, err := io.ReadAll(r.Body)
  if err != nil {
    fmt.Printf("Error reading request body: %v\n", err)
    middleware.SendError(w, "Error reading request body", http.StatusBadRequest)
    return
  }
  
  fmt.Printf("Request body: %s\n", string(bodyBytes))
  
  // Создаем новый Reader для декодирования
  bodyReader := bytes.NewReader(bodyBytes)
  if err := json.NewDecoder(bodyReader).Decode(&requestData); err != nil {
    // Логируем ошибку декодирования
    fmt.Printf("Error decoding request body: %v\n", err)
    middleware.SendError(w, "Invalid request body", http.StatusBadRequest)
    return
  }
  
  // Конвертируем ID в число
  var collectionID int
  switch v := requestData.ID.(type) {
  case float64:
    collectionID = int(v)
  case string:
    if collectionID, err = strconv.Atoi(v); err != nil {
      middleware.SendError(w, "Invalid collection ID format", http.StatusBadRequest)
      return
    }
  case int:
    collectionID = v
  default:
    middleware.SendError(w, "Invalid collection ID type", http.StatusBadRequest)
    return
  }

  // Проверяем, что ID передан
  if requestData.ID == 0 {
    middleware.SendError(w, "Collection ID is required", http.StatusBadRequest)
    return
  }

  // Получаем коллекцию из базы данных
  collection, err := database.GetComparisonCollectionByID(collectionID, userID)
  if err != nil {
    middleware.SendError(w, "Failed to get comparison collection: "+err.Error(), http.StatusInternalServerError)
    return
  }

  // Проверяем, что коллекция найдена
  if collection == nil {
    middleware.SendError(w, "Comparison collection not found", http.StatusNotFound)
    return
  }

  // Отправляем ответ
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(collection)
}