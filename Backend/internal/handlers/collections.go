package handlers

import (
  "encoding/json"
  "net/http"

  "backend/internal/database"
  "backend/internal/middleware"
  "backend/internal/models"
)

// CreateCollectionHandler создает новую подборку
func CreateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем информацию о пользователе из контекста
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		middleware.SendError(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// Декодируем тело запроса
	var collection models.Collection
	if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
		middleware.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Устанавливаем ID пользователя
	collection.UserID = userID

	// Создаем подборку в базе данных
	if err := database.CreateCollection(&collection); err != nil {
		middleware.SendError(w, "Failed to create collection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(collection)
}

// GetCollectionsHandler получает все подборки пользователя
func GetCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем информацию о пользователе из контекста
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		middleware.SendError(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// Получаем подборки из базы данных
	collections, err := database.GetCollectionsByUserID(userID)
	if err != nil {
		middleware.SendError(w, "Failed to get collections: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(collections)
}

// UpdateCollectionHandler обновляет подборку
func UpdateCollectionHandler(w http.ResponseWriter, r *http.Request) {
  // Получаем информацию о пользователе из контекста
  ctx := r.Context()
  userID, ok := ctx.Value("user_id").(int)
  if !ok {
    middleware.SendError(w, "User not found in context", http.StatusInternalServerError)
    return
  }

  // Декодируем тело запроса
  var requestData struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
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

  // Создаем объект подборки
  collection := models.Collection{
    ID:   requestData.ID,
    UserID: userID,
    Name: requestData.Name,
  }

  // Обновляем подборку в базе данных
  if err := database.UpdateCollection(&collection); err != nil {
    if err == database.ErrNoRows {
      middleware.SendError(w, "Collection not found", http.StatusNotFound)
      return
    }
    middleware.SendError(w, "Failed to update collection: "+err.Error(), http.StatusInternalServerError)
    return
  }

  // Отправляем ответ
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(collection)
}

// DeleteCollectionHandler удаляет подборку
func DeleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
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

  // Удаляем подборку из базы данных
  if err := database.DeleteCollection(requestData.ID, userID); err != nil {
    if err == database.ErrNoRows {
      middleware.SendError(w, "Collection not found", http.StatusNotFound)
      return
    }
    middleware.SendError(w, "Failed to delete collection: "+err.Error(), http.StatusInternalServerError)
    return
  }

  // Отправляем ответ
  w.WriteHeader(http.StatusNoContent)
}

// GetCollectionByIDHandler получает подборку по ID
func GetCollectionByIDHandler(w http.ResponseWriter, r *http.Request) {
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

  // Получаем подборку из базы данных
  collection, err := database.GetCollectionByID(requestData.ID, userID)
  if err != nil {
    middleware.SendError(w, "Failed to get collection: "+err.Error(), http.StatusInternalServerError)
    return
  }

  // Проверяем, что подборка найдена
  if collection == nil {
    middleware.SendError(w, "Collection not found", http.StatusNotFound)
    return
  }

  // Отправляем ответ
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(collection)
}