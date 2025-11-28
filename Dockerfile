# Стадия сборки
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bot .

# Стадия запуска
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bot .

# Добавляем healthcheck для мониторинга
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ps aux | grep bot | grep -v grep || exit 1

ENV BOT_TOKEN=""
ENV AI_API_KEY=""
ENV DEBUG="false"
ENV CONFIG_EMAIL=""
ENV BOT_NAME=""

CMD ["./bot"]