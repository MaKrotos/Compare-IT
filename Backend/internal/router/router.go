package router

import (
  "context"
  "net/http"
  "regexp"
  "strings"
  "sync"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type route struct {
  pattern     *regexp.Regexp
  paramNames  []string
  handler     http.HandlerFunc
}

type Router struct {
  middlewares []Middleware
  routes      map[string][]route
  mu          sync.RWMutex
}

func New() *Router {
  return &Router{
    routes: make(map[string][]route),
  }
}

func (r *Router) Use(middleware Middleware) {
  r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) HandleFunc(method, pattern string, handler http.HandlerFunc) {
  r.mu.Lock()
  defer r.mu.Unlock()

  // Извлекаем имена параметров
  paramNames := []string{}
  paramPattern := regexp.MustCompile(":([a-zA-Z0-9_]+)")
  matches := paramPattern.FindAllStringSubmatch(pattern, -1)
  for _, match := range matches {
    if len(match) > 1 {
      paramNames = append(paramNames, match[1])
    }
  }

  // Преобразуем паттерн в регулярное выражение
  regexPattern := "^" + strings.ReplaceAll(regexp.QuoteMeta(pattern), "\\:", "([^/]+)") + "$"
  re := regexp.MustCompile(regexPattern)

  // Применяем middleware к обработчику
  wrappedHandler := handler
  for i := len(r.middlewares) - 1; i >= 0; i-- {
    wrappedHandler = r.middlewares[i](wrappedHandler)
  }

  // Добавляем маршрут
  route := route{
    pattern:    re,
    paramNames: paramNames,
    handler:     wrappedHandler,
  }

  r.routes[method] = append(r.routes[method], route)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  r.mu.RLock()
  defer r.mu.RUnlock()

  if routes, ok := r.routes[req.Method]; ok {
    for _, route := range routes {
      if matches := route.pattern.FindStringSubmatch(req.URL.Path); matches != nil {
        // Создаем контекст с параметрами
        ctx := req.Context()
        // Создаем карту параметров
        params := make(map[string]string)
        for i, paramName := range route.paramNames {
          // i+1 потому что matches[0] - это вся строка
          if i+1 < len(matches) {
            params[paramName] = matches[i+1]
          }
        }
        ctx = context.WithValue(ctx, "url_params", params)
        route.handler(w, req.WithContext(ctx))
        return
      }
    }
  }

  http.NotFound(w, req)
}