# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=keerja-backend
MAIN_PATH=./cmd/api

# Docker parameters
DOCKER_COMPOSE=docker-compose
APP_VERSION?=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: all build clean test coverage run dev docker-up docker-down docker-logs docker-reset help install db-migration-create db-migrate-up db-migrate-down db-migration-status db-migrate-to

help:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║           Keerja Backend - Available Commands                 ║"
	@echo "╠══════════════════════════════════════════════════════════════╣"
	@echo "║ Development:                                                  ║"
	@echo "║   make install      - Download dependencies                   ║"
	@echo "║   make build        - Build application                       ║"
	@echo "║   make run          - Run application                         ║"
	@echo "║   make dev          - Run with hot-reload (requires air)      ║"
	@echo "║   make test         - Run unit tests                          ║"
	@echo "║   make coverage     - Run tests with coverage                 ║"
	@echo "║   make clean        - Clean build files                       ║"
	@echo "╠══════════════════════════════════════════════════════════════╣"
	@echo "║ Docker:                                                       ║"
	@echo "║   make docker-up         - Start infrastructure (db, redis)   ║"
	@echo "║   make docker-dev        - Start with dev tools               ║"
	@echo "║   make docker-app        - Start with API service             ║"
	@echo "║   make docker-full       - Start all services                 ║"
	@echo "║   make docker-down       - Stop all containers                ║"
	@echo "║   make docker-logs       - Show container logs                ║"
	@echo "║   make docker-build      - Build Docker image                 ║"
	@echo "║   make docker-push       - Push to registry                   ║"
	@echo "║   make docker-reset      - Reset database (WARNING!)          ║"
	@echo "╠══════════════════════════════════════════════════════════════╣"
	@echo "║ Database:                                                     ║"
	@echo "║   make db-migrate-up     - Run pending migrations             ║"
	@echo "║   make db-migrate-down   - Rollback one migration             ║"
	@echo "║   make db-migration-status - Show migration version           ║"
	@echo "║   make db-migration-create name=xxx - Create migration        ║"
	@echo "║   make seed              - Run database seeders               ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"

## install: Download semua dependencies
install:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## build: Build aplikasi
build:
	@echo "Building application..."
	$(GOBUILD) -o bin/$(BINARY_NAME) $(MAIN_PATH)

## run: Run aplikasi
run:
	@echo "Running application..."
	$(GOCMD) run $(MAIN_PATH)/main.go

## dev: Run aplikasi dengan auto-reload (butuh cosmtrek/air)
dev:
	@echo "Running in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Error: 'air' is not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

## test: Run unit tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

## coverage: Run tests dengan coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## clean: Clean build files
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -f coverage.out coverage.html

## docker-up: Start infrastructure containers (postgres, redis)
docker-up:
	@echo "Starting infrastructure containers..."
	$(DOCKER_COMPOSE) up -d postgres redis
	@echo "✅ Infrastructure ready!"
	@echo "   PostgreSQL: localhost:5434"
	@echo "   Redis:      localhost:6379"

## docker-dev: Start with development tools (mailhog, adminer, hot-reload api)
docker-dev:
	@echo "Starting development environment..."
	$(DOCKER_COMPOSE) --profile dev up -d
	@echo "✅ Development environment ready!"
	@echo "   API (hot-reload): localhost:8080"
	@echo "   PostgreSQL:       localhost:5434"
	@echo "   Redis:            localhost:6379"
	@echo "   MailHog UI:       localhost:8025"
	@echo "   Adminer:          localhost:8081"

## docker-app: Start with production API
docker-app:
	@echo "Starting with production API..."
	$(DOCKER_COMPOSE) --profile app up -d
	@echo "Application ready!"
	@echo "   API:        localhost:8080"
	@echo "   Health:     localhost:8080/health"

## docker-full: Start all services
docker-full:
	@echo "Starting all services..."
	$(DOCKER_COMPOSE) --profile full up -d
	@echo "All services ready!"

## docker-down: Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) --profile full down

## docker-logs: Show Docker logs
docker-logs:
	@echo "Showing Docker logs..."
	$(DOCKER_COMPOSE) logs -f

## docker-reset: Reset database (WARNING: deletes data)
docker-reset:
	@echo "WARNING: This will delete all database data!"
	@read -p "Are you sure? [y/N] " ans && [ $${ans:-N} = y ]
	$(DOCKER_COMPOSE) --profile full down -v
	$(DOCKER_COMPOSE) up -d postgres redis
	@echo "Database reset complete"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	$(DOCKER_COMPOSE) build api \
		--build-arg APP_VERSION=$(APP_VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT)
	@echo "Image built: keerja-api:$(APP_VERSION)"

## docker-build-dev: Build development Docker image
docker-build-dev:
	@echo "Building development Docker image..."
	$(DOCKER_COMPOSE) build api-dev

## docker-push: Push image to registry
docker-push:
	@if [ -z "$(REGISTRY)" ]; then \
		echo "Error: REGISTRY is required. Usage: make docker-push REGISTRY=your-registry.com"; \
		exit 1; \
	fi
	docker tag keerja-api:latest $(REGISTRY)/keerja-api:$(APP_VERSION)
	docker tag keerja-api:latest $(REGISTRY)/keerja-api:latest
	docker push $(REGISTRY)/keerja-api:$(APP_VERSION)
	docker push $(REGISTRY)/keerja-api:latest
	@echo "Image pushed to $(REGISTRY)"

## docker-restart: Restart Docker containers
docker-restart:
	@echo "Restarting Docker containers..."
	$(DOCKER_COMPOSE) restart

## docker-ps: Show running containers
docker-ps:
	@$(DOCKER_COMPOSE) ps

## docker-stats: Show container stats
docker-stats:
	@docker stats --no-stream $(shell $(DOCKER_COMPOSE) ps -q)

## db-migration-create: Create a new migration
db-migration-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: 'name' parameter is required. Usage: make db-migration-create name=migration_name"; \
		exit 1; \
	fi
	@if ! command -v migrate > /dev/null; then \
		echo "Installing migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir database/migrations -seq $(name)

## db-migrate-up: Run all pending migrations
db-migrate-up:
	@if [ -z "$$DATABASE_URL" ]; then \
		echo "Error: DATABASE_URL environment variable is not set"; \
		echo "Set it with: export DATABASE_URL=\"postgresql://user:password@localhost:5434/dbname?sslmode=disable\""; \
		exit 1; \
	fi
	@if ! command -v migrate > /dev/null; then \
		echo "Installing migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@echo "Running migrations..."
	migrate -path database/migrations -database "$$DATABASE_URL" up
	@echo "Migrations completed successfully"

## db-migrate-up-docker: Apply migrations via docker (using migrator role)
db-migrate-up-docker:
	@echo "Running migrations via Docker..."
	$(DOCKER_COMPOSE) exec migrate migrate -path /migrations -database "postgresql://kustan:$${DB_MIGRATOR_PASSWORD:-kustan_dev_pass_123}@postgres:5432/keerja?sslmode=disable" up
	@echo "Docker migrations completed successfully"

## db-migrate-down: Rollback one migration step
db-migrate-down:
	@if [ -z "$$DATABASE_URL" ]; then \
		echo "Error: DATABASE_URL environment variable is not set"; \
		echo "Set it with: export DATABASE_URL=\"postgresql://user:password@localhost:5434/dbname?sslmode=disable\""; \
		exit 1; \
	fi
	@if ! command -v migrate > /dev/null; then \
		echo "Installing migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@echo "Rolling back one migration..."
	migrate -path database/migrations -database "$$DATABASE_URL" down -steps 1
	@echo "Rollback completed successfully"

## db-migration-status: Show current migration version
db-migration-status:
	@if [ -z "$$DATABASE_URL" ]; then \
		echo "Error: DATABASE_URL environment variable is not set"; \
		echo "Set it with: export DATABASE_URL=\"postgresql://user:password@localhost:5432/dbname?sslmode=disable\""; \
		exit 1; \
	fi
	@if ! command -v migrate > /dev/null; then \
		echo "Installing migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@echo "Current migration version:"
	migrate -path database/migrations -database "$$DATABASE_URL" version

## db-migrate-to: Migrate to specific version
db-migrate-to:
	@if [ -z "$(version)" ]; then \
		echo "Error: 'version' parameter is required. Usage: make db-migrate-to version=2"; \
		exit 1; \
	fi
	@if [ -z "$$DATABASE_URL" ]; then \
		echo "Error: DATABASE_URL environment variable is not set"; \
		echo "Set it with: export DATABASE_URL=\"postgresql://user:password@localhost:5432/dbname?sslmode=disable\""; \
		exit 1; \
	fi
	@if ! command -v migrate > /dev/null; then \
		echo "Installing migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@echo "Migrating to version: $(version)"
	migrate -path database/migrations -database "$$DATABASE_URL" goto $(version)

## seed: Run database seeders
seed:
	@echo "Running database seeders..."
	$(GOCMD) run ./cmd/seeder/main.go

.DEFAULT_GOAL := help
