package handlers

import (
	"backend/internal/router"
)

// RegisterRoutes регистрирует все маршруты приложения
func RegisterRoutes(r *router.Router) {
	r.HandleFunc("GET", "/preview", PreviewHandler)
}