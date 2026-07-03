# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -o /app/bot \
    ./cmd/demobot

# ---------- Runtime stage ----------
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/bot .

USER app

CMD ["./bot"]