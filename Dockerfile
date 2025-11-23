FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Синхронизируем зависимости и собираем приложение
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/main .

# Копируем Swagger документацию
COPY --from=builder /app/docs ./docs

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
