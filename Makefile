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

.PHONY: all build clean test coverage run dev docker-up docker-down docker-logs help

## help: Menampilkan bantuan
help:
	@echo "Available commands:"
	@echo "  make install      - Download dependencies"
	@echo "  make build        - Build aplikasi"
	@echo "  make run          - Run aplikasi"
	@echo "  make dev          - Run aplikasi dengan auto-reload (butuh air)"
	@echo "  make test         - Run unit tests"
	@echo "  make coverage     - Run tests dengan coverage"
	@echo "  make clean        - Clean build files"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make docker-logs  - Show Docker logs"
	@echo "  make migrate-up   - Run database migrations (jika ada)"
	@echo "  make migrate-down - Rollback database migrations (jika ada)"
	@echo "  make seed         - Run database seeders"

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

## docker-up: Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d

## docker-down: Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

## docker-logs: Show Docker logs
docker-logs:
	@echo "Showing Docker logs..."
	$(DOCKER_COMPOSE) logs -f

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	$(DOCKER_COMPOSE) build

## docker-restart: Restart Docker containers
docker-restart:
	@echo "Restarting Docker containers..."
	$(DOCKER_COMPOSE) restart

## migrate-up: Run database migrations (placeholder)
migrate-up:
	@echo "Running migrations..."
	@echo "Migrations not implemented yet. Use your migration tool here."

## migrate-down: Rollback database migrations (placeholder)
migrate-down:
	@echo "Rolling back migrations..."
	@echo "Migrations not implemented yet. Use your migration tool here."

## seed: Run database seeders
seed:
	@echo "Running database seeders..."
	$(GOCMD) run ./cmd/seeder/main.go

.DEFAULT_GOAL := help
