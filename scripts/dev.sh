#!/bin/bash

# Keerja Backend - Development Helper Script

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Install dependencies
install_deps() {
    print_info "Installing Go dependencies..."
    
    # Core dependencies (should already exist)
    go get github.com/gofiber/fiber/v2
    go get gorm.io/gorm
    go get gorm.io/driver/postgres
    go get github.com/joho/godotenv
    go get github.com/go-playground/validator/v10
    go get github.com/sirupsen/logrus
    go get github.com/google/uuid
    go get github.com/gosimple/slug
    go get golang.org/x/crypto
    
    # New required dependencies
    go get github.com/golang-jwt/jwt/v5
    go get github.com/go-mail/mail/v2
    
    # Optional: AWS S3
    # go get github.com/aws/aws-sdk-go-v2
    # go get github.com/aws/aws-sdk-go-v2/service/s3
    
    # Optional: Cloudinary
    # go get github.com/cloudinary/cloudinary-go/v2
    
    # Optional: Redis
    # go get github.com/redis/go-redis/v9
    
    # Optional: Swagger
    # go get github.com/swaggo/swag/cmd/swag
    # go get github.com/gofiber/swagger
    
    go mod tidy
    
    print_success "Dependencies installed"
}

# Build the application
build() {
    print_info "Building application..."
    go build -o bin/keerja-api cmd/api/main.go
    print_success "Build completed: bin/keerja-api"
}

# Run the application
run() {
    print_info "Running application..."
    go run cmd/api/main.go
}

# Run tests
test() {
    print_info "Running tests..."
    go test ./... -v
    print_success "Tests completed"
}

# Run tests with coverage
test_coverage() {
    print_info "Running tests with coverage..."
    go test ./... -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    print_success "Coverage report generated: coverage.html"
}

# Format code
fmt() {
    print_info "Formatting code..."
    go fmt ./...
    print_success "Code formatted"
}

# Run linter
lint() {
    print_info "Running linter..."
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run
        print_success "Linting completed"
    else
        print_error "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    fi
}

# Clean build artifacts
clean() {
    print_info "Cleaning build artifacts..."
    rm -rf bin/
    rm -f coverage.out coverage.html
    print_success "Clean completed"
}

# Database migrations (placeholder)
migrate_up() {
    print_info "Running database migrations..."
    # TODO: Implement migration logic
    print_info "Migration logic not yet implemented"
}

migrate_down() {
    print_info "Rolling back database migrations..."
    # TODO: Implement migration rollback
    print_info "Migration rollback logic not yet implemented"
}

# Database seeder (placeholder)
seed() {
    print_info "Seeding database..."
    # TODO: Implement seeder logic
    print_info "Seeder logic not yet implemented"
}

# Development server with hot reload (requires air)
dev() {
    print_info "Starting development server with hot reload..."
    if command -v air &> /dev/null; then
        air
    else
        print_error "air not installed. Install with: go install github.com/cosmtrek/air@latest"
        print_info "Running without hot reload..."
        run
    fi
}

# Show help
show_help() {
    echo "Keerja Backend - Development Helper"
    echo ""
    echo "Usage: ./scripts/dev.sh [command]"
    echo ""
    echo "Commands:"
    echo "  install        Install all dependencies"
    echo "  build          Build the application"
    echo "  run            Run the application"
    echo "  dev            Run with hot reload (requires air)"
    echo "  test           Run all tests"
    echo "  test-coverage  Run tests with coverage report"
    echo "  fmt            Format code"
    echo "  lint           Run linter"
    echo "  clean          Clean build artifacts"
    echo "  migrate-up     Run database migrations"
    echo "  migrate-down   Rollback database migrations"
    echo "  seed           Seed database with initial data"
    echo "  help           Show this help message"
    echo ""
}

# Main script
case "$1" in
    install)
        install_deps
        ;;
    build)
        build
        ;;
    run)
        run
        ;;
    dev)
        dev
        ;;
    test)
        test
        ;;
    test-coverage)
        test_coverage
        ;;
    fmt)
        fmt
        ;;
    lint)
        lint
        ;;
    clean)
        clean
        ;;
    migrate-up)
        migrate_up
        ;;
    migrate-down)
        migrate_down
        ;;
    seed)
        seed
        ;;
    help|"")
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
