# Keerja Backend API

> Job portal backend API built with Go Fiber, GORM, PostgreSQL using Clean Architecture

[![CI](https://github.com/rasfaxo/keerja-backend/actions/workflows/ci.yml/badge.svg)](https://github.com/rasfaxo/keerja-backend/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org)
[![Fiber Version](https://img.shields.io/badge/Fiber-v2.52-blue.svg)](https://gofiber.io)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

---

## Table of Contents

- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Development](#development)
- [API Documentation](#api-documentation)
- [Tech Stack](#tech-stack)
- [Testing](#testing)
- [Docker Services](#docker-services)
- [Security](#security)
- [Contributing](#contributing)
- [License](#license)

---

## Architecture

Project ini menggunakan **Clean Architecture** dengan **Domain-Driven Design (DDD)**:

```
keerja-backend/
├── cmd/                  # Application entry points
├── internal/             # Private application code
│   ├── config/          # Configuration management
│   ├── domain/          # Business entities & interfaces
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic layer
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── dto/             # Data transfer objects
│   ├── routes/          # Route definitions
│   └── utils/           # Utility functions
├── pkg/                  # Public reusable packages
│   ├── logger/          # Logging
│   ├── storage/         # File storage
│   └── email/           # Email service
├── database/            # Database management
│   ├── migrations/      # Migration files
│   └── seeders/         # Seed data
├── tests/               # Tests
└── docs/                # Documentation
```

**Principles:**

- Separation of Concerns
- Dependency Inversion
- Single Responsibility
- Interface Segregation
- Repository Pattern
- Service Layer Pattern

---

## Quick Start

### Prerequisites

- **Go** 1.23 or higher
- **Docker** & **Docker Compose** (recommended)
- **Make** (optional, untuk menggunakan Makefile commands)

### Installation

#### Option 1: Using Docker (Recommended)

1. **Clone repository:**

```bash
git clone https://github.com/rasfaxo/keerja-backend.git
cd keerja-backend
```

2. **Setup environment:**

```bash
cp .env.example .env
```

3. **Start Docker containers:**

```bash
# Start infrastructure only (PostgreSQL, Redis)
make docker-up

# Start with development tools (+ MailHog, Adminer)
make docker-dev

# Start all services including API
make docker-full
```

This will start:
| Service | Port | Description |
|---------|------|-------------|
| PostgreSQL 17 | 5434 | Database |
| Redis | 6379 | Cache & Session |
| MinIO | 9000, 9001 | Object Storage |
| MailHog | 1025, 8025 | Email Testing |
| Adminer | 8081 | Database Management UI |

4. **Run database migrations:**

```bash
make db-migrate-up
```

5. **Seed database with initial data:**

```bash
make seed
```

This seeds:

- 34 Provinces, 57 Cities, 42 Districts (Indonesia)
- 36 Industries
- 6 Company Sizes
- 5 Job Types, 3 Work Policies
- 7 Education Levels, 6 Experience Levels
- 30 Job Titles, 49 Skills, 49 Benefits
- 68 Job Categories, 37 Job Subcategories
- 7 Admin Roles + 1 Admin User
- 15 Sample Companies

6. **Run the application:**

```bash
make dev
# atau
make run
```

7. **Access the services:**

| Service       | URL                          | Credentials                    |
| ------------- | ---------------------------- | ------------------------------ |
| API           | http://localhost:8080        | -                              |
| API Health    | http://localhost:8080/health | -                              |
| Adminer       | http://localhost:8081        | See below                      |
| MailHog UI    | http://localhost:8025        | -                              |
| MinIO Console | http://localhost:9001        | minioadmin / minioadmin123     |
| PostgreSQL    | localhost:5434               | bekerja / bekerja_dev_pass_123 |

**Adminer Login:**

- System: PostgreSQL
- Server: postgres
- Username: postgres
- Password: postgres_admin_pass
- Database: keerja

#### Option 2: Manual Setup

1. **Clone & Install dependencies:**

```bash
git clone https://github.com/rasfaxo/keerja-backend.git
cd keerja-backend
make install
# atau
go mod download && go mod tidy
```

2. **Setup PostgreSQL manually:**

```bash
# Create database
createdb -U postgres keerja

# Create user
psql -U postgres -c "CREATE USER bekerja WITH PASSWORD 'bekerja_dev_pass_123';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE keerja TO bekerja;"
```

3. **Setup environment:**

```bash
cp .env.example .env
# Edit DATABASE_URL in .env:
# DATABASE_URL=postgresql://bekerja:bekerja_dev_pass_123@localhost:5432/keerja?sslmode=disable
```

4. **Run migrations & seeders:**

```bash
make db-migrate-up
make seed
```

5. **Run application:**

```bash
make run
```

### Default Admin Account

After seeding, you can login with:

- **Email:** admin@keerja.com
- **Password:** admin123

### Environment Variables

Key environment variables in `.env`:

```env
# Application
APP_ENV=development
APP_PORT=8080

# Database (Docker uses port 5434)
DATABASE_URL=postgresql://bekerja:bekerja_dev_pass_123@localhost:5434/keerja?sslmode=disable

# JWT
JWT_SECRET=keerja_jwt_secret_key_minimum_32_characters_long_for_security
JWT_EXPIRE_HOURS=24
JWT_REFRESH_EXPIRE_DAYS=7

# Redis
REDIS_URL=redis://localhost:6379/0

# Firebase (optional)
FCM_ENABLED=false
FCM_CREDENTIALS_FILE=./config/firebase-service-account.json
```

API akan berjalan di `http://localhost:8080`

---

## Development

### Available Commands

```bash
# Menggunakan Makefile
make help           # Show all available commands
make install        # Install dependencies
make build          # Build aplikasi
make run            # Run aplikasi
make dev            # Run aplikasi dengan auto-reload (butuh air)
make test           # Run unit tests
make coverage       # Run tests dengan coverage
make clean          # Clean build artifacts

# Docker Infrastructure
make docker-up      # Start infrastructure (postgres, redis)
make docker-dev     # Start with dev tools (mailhog, adminer)
make docker-app     # Start with production API
make docker-full    # Start all services
make docker-down    # Stop all containers
make docker-logs    # Show Docker logs
make docker-build   # Build Docker image
make docker-reset   # Reset database (WARNING: deletes data)

# Database Commands
make db-migrate-up      # Run all pending migrations
make db-migrate-down    # Rollback one migration step
make db-migration-status # Show current migration version
make seed               # Run database seeders
```

### Development with Hot Reload

Install Air untuk hot reload:

```bash
go install github.com/cosmtrek/air@latest
```

Kemudian:

```bash
make dev
```

### Running Tests

```bash
# Run all tests
make test

# Run tests dengan coverage
make coverage

# Run specific package tests
go test ./internal/domain/user/... -v
```

### Linting and Formatting

To ensure code quality and proper formatting, you can use the following commands:

- **Lint**: Run the linter to check for issues in the code.

  ```bash
  make lint
  ```

- **Format**: Automatically format the code.
  ```bash
  make fmt
  ```

These commands use `golangci-lint` for linting and `gofmt` for formatting.

---

## API Documentation

### Health Check Endpoints

```
GET    /health                   General health status
GET    /health/live              Kubernetes liveness probe
GET    /health/ready             Kubernetes readiness probe (checks DB & Redis)
GET    /health/system            System information (admin)
```

### Authentication Endpoints

```
POST   /api/v1/auth/register           Register new user
POST   /api/v1/auth/login              Login (returns JWT tokens)
POST   /api/v1/auth/verify-email       Verify email with OTP
POST   /api/v1/auth/forgot-password    Request password reset OTP
POST   /api/v1/auth/reset-password     Reset password with OTP
POST   /api/v1/auth/refresh-token      Refresh access token
POST   /api/v1/auth/logout             Logout (invalidate tokens)
GET    /api/v1/auth/oauth/google       Google OAuth login
POST   /api/v1/auth/oauth/google/mobile   Mobile Google OAuth (PKCE)
```

### User Endpoints

```
GET    /api/v1/users/me                   Get current user profile
PUT    /api/v1/users/me                   Update profile
POST   /api/v1/users/me/education         Add education
PUT    /api/v1/users/me/education/:id     Update education
DELETE /api/v1/users/me/education/:id     Delete education
POST   /api/v1/users/me/experience        Add experience
POST   /api/v1/users/me/skills            Add skills
```

### Job Endpoints

```
GET    /api/v1/jobs            List jobs (with filters)
GET    /api/v1/jobs/:id        Get job details
POST   /api/v1/jobs/search     Advanced job search
POST   /api/v1/jobs            Create job (employer only)
PUT    /api/v1/jobs/:id        Update job (employer only)
DELETE /api/v1/jobs/:id        Delete job (employer only)
```

### Master Data Endpoints

```
GET    /api/v1/master/provinces         List provinces
GET    /api/v1/master/cities            List cities
GET    /api/v1/master/districts         List districts
GET    /api/v1/master/industries        List industries
GET    /api/v1/master/company-sizes     List company sizes
GET    /api/v1/master/job-types         List job types
GET    /api/v1/master/work-policies     List work policies
GET    /api/v1/master/education-levels  List education levels
GET    /api/v1/master/experience-levels List experience levels
GET    /api/v1/master/skills            List skills
GET    /api/v1/master/benefits          List benefits
GET    /api/v1/master/job-categories    List job categories
```

### Push Notification Endpoints (FCM)

```
POST   /api/v1/device-tokens           Register device token
GET    /api/v1/device-tokens           Get user's devices
DELETE /api/v1/device-tokens/:token    Unregister device
POST   /api/v1/push/send/user/:id      Send notification to user
POST   /api/v1/push/send/batch         Send to multiple users
POST   /api/v1/push/send/topic         Send to topic subscribers
```

**Authentication**: All protected endpoints require `Authorization: Bearer <token>` header.

---

## Tech Stack

| Category          | Technology                  |
| ----------------- | --------------------------- |
| Language          | Go 1.24                     |
| Framework         | Fiber v2.52                 |
| ORM               | GORM v1.25                  |
| Database          | PostgreSQL 17               |
| Cache             | Redis                       |
| Auth              | JWT (golang-jwt/jwt/v5)     |
| Password          | Bcrypt                      |
| Validation        | go-playground/validator/v10 |
| Logging           | Logrus                      |
| Email             | gomail.v2 + MailHog (dev)   |
| Storage           | MinIO (S3-compatible)       |
| Push Notification | Firebase Cloud Messaging    |
| Container         | Docker & Docker Compose     |

---

## Testing

```bash
# Run all tests
make test

# Run specific package
go test ./internal/service/... -v

# Run with coverage
make coverage

# View coverage in browser
go tool cover -html=coverage.out
```

---

## Docker Services

```bash
# Start infrastructure only
make docker-up

# Start with development tools
make docker-dev

# Start all services
make docker-full

# Stop all services
make docker-down

# View logs
make docker-logs

# Build production image
make docker-build

# Reset database (WARNING: deletes all data)
make docker-reset
```

| Service    | Port       | URL                   |
| ---------- | ---------- | --------------------- |
| API        | 8080       | http://localhost:8080 |
| PostgreSQL | 5434       | localhost:5434        |
| Redis      | 6379       | localhost:6379        |
| MinIO      | 9000, 9001 | http://localhost:9001 |
| MailHog    | 1025, 8025 | http://localhost:8025 |
| Adminer    | 8081       | http://localhost:8081 |

---

## Security

- Password hashing dengan `bcrypt`
- JWT Bearer Token authentication
- Input validation dengan `go-playground/validator`
- SQL injection prevention dengan GORM prepared statements
- CORS middleware
- Rate limiting
- GitHub Actions CI with security scanning (Trivy)
- Docker multi-stage builds with non-root user

---

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

This project is licensed under the MIT License.

---

## Authors

- **rasfaxo** - [GitHub](https://github.com/rasfaxo)

---

**Built with Go**
