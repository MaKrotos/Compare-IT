package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		
		// Вызываем следующий обработчик
		next(w, r)
		
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}
}