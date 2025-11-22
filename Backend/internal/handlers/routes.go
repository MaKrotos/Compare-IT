package handlers

import (
	"backend/internal/middleware"
	"backend/internal/router"
)

// RegisterRoutes регистрирует все маршруты приложения
func RegisterRoutes(r *router.Router) {
	// Публичные маршруты
	r.HandleFunc("GET", "/preview", PreviewHandler)
	
	// Аутентификация
	r.HandleFunc("POST", "/auth/telegram", TelegramAuthHandler)
	
	// Защищенные маршруты
	r.HandleFunc("GET", "/profile", middleware.JWTAuthMiddleware(ProfileHandler))
	
	// Маршруты для работы с коллекциями сравнений
	r.HandleFunc("POST", "/comparisons", middleware.JWTAuthMiddleware(CreateComparisonCollectionHandler))
	r.HandleFunc("GET", "/comparisons", middleware.JWTAuthMiddleware(GetComparisonCollectionsHandler))
	r.HandleFunc("POST", "/comparisons/get", middleware.JWTAuthMiddleware(GetComparisonCollectionByIdHandler))
	r.HandleFunc("POST", "/comparisons/update", middleware.JWTAuthMiddleware(UpdateComparisonCollectionHandler))
	r.HandleFunc("POST", "/comparisons/delete", middleware.JWTAuthMiddleware(DeleteComparisonCollectionHandler))
	
	// Маршруты для работы с подборками
	r.HandleFunc("POST", "/collections", middleware.JWTAuthMiddleware(CreateCollectionHandler))
	r.HandleFunc("GET", "/collections", middleware.JWTAuthMiddleware(GetCollectionsHandler))
	r.HandleFunc("POST", "/collections/update", middleware.JWTAuthMiddleware(UpdateCollectionHandler))
	r.HandleFunc("POST", "/collections/delete", middleware.JWTAuthMiddleware(DeleteCollectionHandler))
	r.HandleFunc("POST", "/collections/get", middleware.JWTAuthMiddleware(GetCollectionByIDHandler))
	
	// Маршруты для работы с публичными ссылками
	r.HandleFunc("POST", "/collections/generate-public-link", middleware.JWTAuthMiddleware(GeneratePublicLinkHandler))
	r.HandleFunc("POST", "/collections/remove-public-link", middleware.JWTAuthMiddleware(RemovePublicLinkHandler))
	r.HandleFunc("POST", "/collections/public", GetPublicCollectionHandler)
}