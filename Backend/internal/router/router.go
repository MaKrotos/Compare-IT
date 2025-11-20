package router

import (
	"net/http"
	"sync"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Router struct {
	middlewares []Middleware
	routes      map[string]map[string]http.HandlerFunc
	mu          sync.RWMutex
}

func New() *Router {
	return &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

func (r *Router) Use(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) HandleFunc(method, pattern string, handler http.HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.routes[method] == nil {
		r.routes[method] = make(map[string]http.HandlerFunc)
	}

	// Применяем middleware к обработчику
	wrappedHandler := handler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		wrappedHandler = r.middlewares[i](wrappedHandler)
	}

	r.routes[method][pattern] = wrappedHandler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if handlers, ok := r.routes[req.Method]; ok {
		if handler, ok := handlers[req.URL.Path]; ok {
			handler(w, req)
			return
		}
	}

	http.NotFound(w, req)
}