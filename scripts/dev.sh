#!/bin/bash
set -e

echo "🚀 Запуск DemoBot (dev mode)"

if [ ! -f .env ]; then
  echo "⚠️  .env не найден"
  cp .env.example .env
  echo "✏️  Заполни .env и перезапусти"
  exit 1
fi

set -a
source .env
set +a

echo "🐘 Проверяем PostgreSQL..."
if ! nc -z localhost 5432; then
  echo "❌ PostgreSQL не запущен на 5432"
  exit 1
fi

echo "📦 Применяем миграции..."
goose -dir build/app/migrations postgres \
"postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up

echo "🤖 Запуск бота..."
go run ./cmd/demobot/main.go
