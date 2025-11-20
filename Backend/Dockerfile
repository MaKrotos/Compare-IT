FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

# Добавляем команду go mod tidy для автоматического управления зависимостями
RUN go mod tidy

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/go.mod .

EXPOSE 8080

CMD ["./main"]