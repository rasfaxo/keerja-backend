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

.PHONY: all build clean test coverage run dev docker-up docker-down docker-logs docker-reset help install db-migration-create db-migrate-up db-migrate-down db-migration-status db-migrate-to lint fmt

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
	@echo "║   make lint         - Run linter                              ║"
	@echo "║   make fmt          - Format code                            ║"
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

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

## fmt: Format code
fmt:
	@echo "Formatting code..."
	gofmt -w .

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

# ============================================================
# VPS Multi-Environment Commands
# ============================================================
# These commands are for managing STAGING and DEMO environments
# on the VPS (145.79.8.227)
#
# Prerequisites:
#   - SSH access to VPS configured
#   - .env.staging and .env.demo files present on VPS
#   - docker-compose.vps.yml synced to VPS
# ============================================================

# VPS Configuration
VPS_HOST ?= 145.79.8.227
VPS_USER ?= root
VPS_PATH ?= /opt/keerja
VPS_COMPOSE_FILE = docker-compose.vps.yml

.PHONY: vps-help vps-ssh vps-status vps-logs-staging vps-logs-demo vps-restart-staging vps-restart-demo vps-deploy-staging vps-deploy-demo vps-backup-staging vps-backup-demo vps-health

## vps-help: Show VPS-related commands
vps-help:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║           Keerja Backend - VPS Commands                       ║"
	@echo "╠══════════════════════════════════════════════════════════════╣"
	@echo "║ Connection:                                                   ║"
	@echo "║   make vps-ssh           - SSH into VPS                       ║"
	@echo "║   make vps-status        - Show container status              ║"
	@echo "║   make vps-health        - Run health check                   ║"
	@echo "╠══════════════════════════════════════════════════════════════╣"
	@echo "║ STAGING (http://staging-api.145.79.8.227.nip.io):             ║"
	@echo "║   make vps-deploy-staging  - Deploy to STAGING                ║"
	@echo "║   make vps-logs-staging    - View STAGING logs                ║"
	@echo "║   make vps-restart-staging - Restart STAGING                  ║"
	@echo "║   make vps-backup-staging  - Backup STAGING database          ║"
	@echo "╠══════════════════════════════════════════════════════════════╣"
	@echo "║ DEMO (http://demo-api.145.79.8.227.nip.io):                   ║"
	@echo "║   make vps-deploy-demo     - Deploy to DEMO                   ║"
	@echo "║   make vps-logs-demo       - View DEMO logs                   ║"
	@echo "║   make vps-restart-demo    - Restart DEMO                     ║"
	@echo "║   make vps-backup-demo     - Backup DEMO database             ║"
	@echo "╠══════════════════════════════════════════════════════════════╣"
	@echo "║ Infrastructure:                                               ║"
	@echo "║   make vps-infra-start   - Start PostgreSQL & Redis           ║"
	@echo "║   make vps-infra-stop    - Stop all VPS containers            ║"
	@echo "║   make vps-migrate-staging - Run STAGING migrations           ║"
	@echo "║   make vps-migrate-demo    - Run DEMO migrations              ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"

## vps-ssh: SSH into VPS
vps-ssh:
	@echo "Connecting to VPS..."
	ssh $(VPS_USER)@$(VPS_HOST)

## vps-status: Show VPS container status
vps-status:
	@echo "Checking VPS container status..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && docker-compose -f $(VPS_COMPOSE_FILE) ps"

## vps-health: Run health check on VPS
vps-health:
	@echo "Running health check on VPS..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && chmod +x scripts/health-check.sh && ./scripts/health-check.sh"

## vps-infra-start: Start infrastructure services on VPS
vps-infra-start:
	@echo "Starting infrastructure on VPS..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && docker-compose -f $(VPS_COMPOSE_FILE) up -d postgres redis"
	@echo "Infrastructure started. Waiting for services..."
	@sleep 10
	@ssh $(VPS_USER)@$(VPS_HOST) "docker exec keerja-vps-postgres pg_isready -U postgres || echo 'PostgreSQL not ready yet'"

## vps-infra-stop: Stop all VPS containers
vps-infra-stop:
	@echo "Stopping all VPS containers..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && docker-compose -f $(VPS_COMPOSE_FILE) --profile staging --profile demo down"

## vps-deploy-staging: Deploy STAGING environment
vps-deploy-staging:
	@echo "Deploying to STAGING..."
	@echo "Step 1: Syncing code..."
	rsync -avz --delete \
		--exclude '.git' \
		--exclude 'node_modules' \
		--exclude '*.log' \
		--exclude '.env' \
		--exclude '.env.staging' \
		--exclude '.env.demo' \
		--exclude 'uploads/*' \
		-e ssh ./ $(VPS_USER)@$(VPS_HOST):$(VPS_PATH)/
	@echo "Step 2: Building and deploying..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && \
		docker-compose -f $(VPS_COMPOSE_FILE) build api-staging && \
		docker-compose -f $(VPS_COMPOSE_FILE) --profile staging up -d api-staging"
	@echo "Step 3: Health check..."
	@sleep 15
	@curl -sf http://$(VPS_HOST):8080/health/live && echo "STAGING deployed successfully!" || echo "Health check failed"

## vps-deploy-demo: Deploy DEMO environment
vps-deploy-demo:
	@echo "Deploying to DEMO..."
	@echo "Step 1: Syncing code..."
	rsync -avz --delete \
		--exclude '.git' \
		--exclude 'node_modules' \
		--exclude '*.log' \
		--exclude '.env' \
		--exclude '.env.staging' \
		--exclude '.env.demo' \
		--exclude 'uploads/*' \
		-e ssh ./ $(VPS_USER)@$(VPS_HOST):$(VPS_PATH)/
	@echo "Step 2: Building and deploying..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && \
		docker-compose -f $(VPS_COMPOSE_FILE) build api-demo && \
		docker-compose -f $(VPS_COMPOSE_FILE) --profile demo up -d api-demo"
	@echo "Step 3: Health check..."
	@sleep 15
	@curl -sf http://$(VPS_HOST):8081/health/live && echo "DEMO deployed successfully!" || echo "Health check failed"

## vps-logs-staging: View STAGING logs
vps-logs-staging:
	@echo "Viewing STAGING logs (Ctrl+C to exit)..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && docker-compose -f $(VPS_COMPOSE_FILE) logs -f api-staging"

## vps-logs-demo: View DEMO logs
vps-logs-demo:
	@echo "Viewing DEMO logs (Ctrl+C to exit)..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && docker-compose -f $(VPS_COMPOSE_FILE) logs -f api-demo"

## vps-restart-staging: Restart STAGING container
vps-restart-staging:
	@echo "Restarting STAGING..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && docker-compose -f $(VPS_COMPOSE_FILE) --profile staging restart api-staging"
	@echo "STAGING restarted"

## vps-restart-demo: Restart DEMO container
vps-restart-demo:
	@echo "Restarting DEMO..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && docker-compose -f $(VPS_COMPOSE_FILE) --profile demo restart api-demo"
	@echo "DEMO restarted"

## vps-migrate-staging: Run migrations on STAGING
vps-migrate-staging:
	@echo "Running STAGING migrations..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && docker-compose -f $(VPS_COMPOSE_FILE) --profile migrate-staging up --build"
	@echo "STAGING migrations complete"

## vps-migrate-demo: Run migrations on DEMO
vps-migrate-demo:
	@echo "Running DEMO migrations..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && docker-compose -f $(VPS_COMPOSE_FILE) --profile migrate-demo up --build"
	@echo "DEMO migrations complete"

## vps-backup-staging: Backup STAGING database
vps-backup-staging:
	@echo "Backing up STAGING database..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && \
		mkdir -p backups && \
		docker exec keerja-vps-postgres pg_dump -U postgres keerja_staging > backups/staging_$$(date +%Y%m%d_%H%M%S).sql"
	@echo "Backup created in $(VPS_PATH)/backups/"

## vps-backup-demo: Backup DEMO database
vps-backup-demo:
	@echo "Backing up DEMO database..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && \
		mkdir -p backups && \
		docker exec keerja-vps-postgres pg_dump -U postgres keerja_demo > backups/demo_$$(date +%Y%m%d_%H%M%S).sql"
	@echo "Backup created in $(VPS_PATH)/backups/"

## vps-seed-staging: Seed STAGING database
vps-seed-staging:
	@echo "Seeding STAGING database..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && docker-compose -f $(VPS_COMPOSE_FILE) --profile seed-staging up --build"
	@echo "STAGING seeding complete"

## vps-seed-demo: Seed DEMO database
vps-seed-demo:
	@echo "Seeding DEMO database..."
	@ssh $(VPS_USER)@$(VPS_HOST) "cd $(VPS_PATH) && docker-compose -f $(VPS_COMPOSE_FILE) --profile seed-demo up --build"
	@echo "DEMO seeding complete"

.DEFAULT_GOAL := help

