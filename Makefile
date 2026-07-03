ifneq (,$(wildcard .env))
include .env
export
endif

COMPOSE := docker compose --env-file .env
BINARY := demobot

.PHONY: help build clean run dev test fmt vet up down restart logs ps health docker-build

help:
	@echo "DemoBot"
	@echo ""
	@echo "Development:"
	@echo "  make dev"
	@echo "  make run"
	@echo ""
	@echo "Docker:"
	@echo "  make up"
	@echo "  make down"
	@echo "  make logs"
	@echo "  make ps"
	@echo "  make health"
	@echo ""
	@echo "Go:"
	@echo "  make build"
	@echo "  make test"

build:
	go build -o $(BINARY) ./cmd/demobot

run:
	go run ./cmd/demobot

dev:
	$(COMPOSE) up -d postgres

	@echo "Waiting for PostgreSQL..."

	@until $(COMPOSE) exec -T postgres \
		pg_isready \
		-U "$(DB_USER)" \
		-d "$(DB_NAME)" >/dev/null 2>&1; do \
		sleep 1; \
	done

	go run ./cmd/demobot

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

restart:
	$(COMPOSE) restart

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

health:
	$(COMPOSE) ps
	@echo ""
	$(COMPOSE) logs --tail=30

docker-build:
	docker build -t demobot:local .

test:
	go test ./... -v

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY)

docker-build:
    docker build -t demobot:local .