# ================================
# STAGE 1: Build
# ================================
FROM golang:1.25-alpine AS builder

# Build arguments
ARG APP_VERSION=1.0.0
ARG BUILD_TIME
ARG GIT_COMMIT

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${APP_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -a -installsuffix cgo \
    -o /app/main ./cmd/main.go

# ================================
# STAGE 2: Production
# ================================
FROM alpine:3.19 AS production

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl

# Set timezone
ENV TZ=Asia/Jakarta

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy config directory if needed (e.g., firebase credentials)
COPY --from=builder /app/config ./config

# Create necessary directories with proper permissions
RUN mkdir -p /app/logs /app/uploads && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health/live || exit 1

# Run the application
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
