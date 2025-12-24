# ================================
# STAGE 1: Build
# ================================
FROM golang:1.25-alpine AS builder

ARG APP_VERSION=1.0.0
ARG BUILD_TIME
ARG GIT_COMMIT

RUN apk add --no-cache git make ca-certificates tzdata
RUN adduser -D -g '' appuser
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# 1. Build Main API Server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${APP_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -a -installsuffix cgo \
    -o /app/main ./cmd/main.go

# 2. Build Migration Tool 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o /app/migrate_tool ./cmd/migrate/main.go

# 3. Build Seeder Tool
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o /app/seeder ./cmd/seeder/main.go

# ================================
# STAGE 2: Production
# ================================
FROM alpine:3.19 AS production

RUN apk --no-cache add ca-certificates tzdata curl
ENV TZ=Asia/Jakarta
RUN adduser -D -g '' appuser
WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrate_tool .
COPY --from=builder /app/seeder .

# Copy config & migrations
COPY --from=builder /app/config ./config
COPY --from=builder /app/database/migrations ./database/migrations

RUN mkdir -p /app/logs /app/uploads && \
    chown -R appuser:appuser /app

USER appuser
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health/live || exit 1

# Default entrypoint 
ENTRYPOINT ["./main"]

# ================================
# STAGE 3: Development (with hot reload)
# ================================
FROM golang:1.25-alpine AS development

# Install development tools
RUN apk add --no-cache git make curl

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code (will be overwritten by volume mount)
COPY . .

# Expose port
EXPOSE 8080

# Run with air for hot reload
CMD ["air", "-c", ".air.toml"]
