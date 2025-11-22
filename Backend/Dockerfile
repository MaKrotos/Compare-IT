FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем только go.mod, если go.sum нет
COPY go.mod ./
RUN go mod download

# Устанавливаем goose для миграций в отдельном слое для кэширования
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

# Добавляем команду go mod tidy для автоматического управления зависимостями
RUN go mod tidy

# Собираем приложение
RUN go build -o main .

FROM alpine:latest

# Устанавливаем зависимости для работы с PostgreSQL
RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

# Копируем бинарный файл приложения
COPY --from=builder /app/main .
COPY --from=builder /app/go.mod .

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Копируем бинарный файл goose
COPY --from=builder /go/bin/goose ./goose

# Делаем скрипт запуска исполняемым
RUN chmod +x ./goose

# Открываем порт
EXPOSE 8080

# Команда запуска приложения с миграциями
CMD ["sh", "-c", "./goose -dir ./migrations postgres $DATABASE_URL up && ./main"]