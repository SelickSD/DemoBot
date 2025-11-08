# Стадия сборки
FROM golang:1.21-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы модулей сначала для кэширования
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение (main.go в корне)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bot .

# Стадия запуска
FROM alpine:latest

# Устанавливаем CA certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Создаем пользователя для безопасности
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Копируем бинарник из стадии сборки
COPY --from=builder /app/bot .

# Меняем владельца файлов
RUN chown -R app:app /app

# Переключаемся на непривилегированного пользователя
USER app

# Команда запуска
CMD ["./bot"]