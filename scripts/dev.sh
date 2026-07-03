#!/bin/bash
set -e

echo "🚀 Starting DemoBot (development)"

if [ ! -f .env ]; then
    echo ".env not found"
    cp .env.example .env
    echo "Please edit .env"
    exit 1
fi

set -a
source .env
set +a

echo "Starting PostgreSQL..."

docker compose up -d postgres

echo "Waiting for PostgreSQL..."

until docker compose exec postgres pg_isready \
    -U "$DB_USER" \
    -d "$DB_NAME" >/dev/null 2>&1
do
    sleep 1
done

echo "Running bot..."

go run ./cmd/demobot