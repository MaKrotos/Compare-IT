# CompareIT - Сервис сравнения

CompareIT - это веб-приложение для сравнения различных продуктов, услуг или идей с помощью интерактивных карточек с плюсами и минусами.

## Описание

Приложение позволяет пользователям создавать сравнительные анализы, добавлять плюсы и минусы для каждого элемента сравнения, а также делиться результатами с другими пользователями. Основной функционал включает предварительный просмотр ссылок с извлечением метаданных (заголовок, описание, изображение).

## Архитектура

Приложение состоит из двух основных компонентов:
- **Frontend**: Vue.js приложение с TailwindCSS
- **Backend**: Go приложение с модульной архитектурой

### Структура Backend

```
Backend/
├── internal/
│   ├── app/          # Основное приложение
│   ├── handlers/     # Обработчики HTTP запросов
│   ├── middleware/   # Middleware функции
│   ├── models/      # Модели данных
│   ├── router/       # Кастомный маршрутизатор
│   └── services/     # Бизнес-логика
└── main.go          # Точка входа
```

## Быстрый запуск

### Требования

- Docker и Docker Compose
- Git

### Запуск приложения

1. Клонируйте репозиторий:
   ```bash
   git clone <url>
   cd CompareIT
   ```

2. Запустите приложение с помощью Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. Откройте приложение в браузере:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080

### Остановка приложения

```bash
docker-compose down
```

## API Endpoints

### Предварительный просмотр URL

- **URL**: `/preview`
- **Метод**: `GET`
- **Параметры**: `url` - URL для предварительного просмотра
- **Ответ**: JSON с метаданными страницы

Пример:
```bash
curl "http://localhost:8080/preview?url=https://example.com"
```

## Разработка

### Backend

Backend написан на Go с использованием модульной архитектуры:

1. Добавление нового endpoint:
   - Создайте обработчик в `internal/handlers/`
   - Зарегистрируйте маршрут в `internal/handlers/routes.go`

2. Добавление нового сервиса:
   - Создайте сервис в `internal/services/`
   - Используйте сервис в обработчиках

3. Добавление middleware:
   - Создайте middleware в `internal/middleware/`
   - Подключите в `internal/app/app.go`

### Frontend

Frontend написан на Vue.js 3 с использованием Composition API:

1. Компоненты UI находятся в `Frontend/src/components/ui/`
2. Компоненты сравнения находятся в `Frontend/src/components/comparison/`
3. Утилиты находятся в `Frontend/src/utils/`

## Сборка и развертывание

### Docker

Приложение использует Docker для контейнеризации:

- `Frontend/Dockerfile` - для frontend приложения
- `Backend/Dockerfile` - для backend приложения
- `docker-compose.yml` - для оркестрации сервисов

### Развёртывание на сервере

1. Установите Docker и Docker Compose на сервере
2. Скопируйте репозиторий на сервер
3. Выполните `docker-compose up -d`

