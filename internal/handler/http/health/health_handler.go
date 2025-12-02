package health

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Status represents the health status of a component
type Status string

const (
	StatusUp   Status = "up"
	StatusDown Status = "down"
)

// ComponentHealth represents health info for a single component
type ComponentHealth struct {
	Status  Status `json:"status"`
	Message string `json:"message,omitempty"`
	Latency string `json:"latency,omitempty"`
}

// HealthResponse represents the full health check response
type HealthResponse struct {
	Status     Status                     `json:"status"`
	Version    string                     `json:"version"`
	Uptime     string                     `json:"uptime"`
	Timestamp  string                     `json:"timestamp"`
	Components map[string]ComponentHealth `json:"components,omitempty"`
}

// ReadinessResponse represents readiness check response
type ReadinessResponse struct {
	Ready      bool                       `json:"ready"`
	Components map[string]ComponentHealth `json:"components"`
}

// LivenessResponse represents liveness check response
type LivenessResponse struct {
	Alive     bool   `json:"alive"`
	Timestamp string `json:"timestamp"`
}

// SystemInfo represents system information
type SystemInfo struct {
	GoVersion     string `json:"go_version"`
	NumCPU        int    `json:"num_cpu"`
	NumGoroutine  int    `json:"num_goroutine"`
	MemAlloc      string `json:"mem_alloc"`
	MemTotalAlloc string `json:"mem_total_alloc"`
	MemSys        string `json:"mem_sys"`
	NumGC         uint32 `json:"num_gc"`
}

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db        *gorm.DB
	redis     *redis.Client
	version   string
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB, redis *redis.Client, version string) *HealthHandler {
	return &HealthHandler{
		db:        db,
		redis:     redis,
		version:   version,
		startTime: time.Now(),
	}
}

// Health returns the overall health status of the application
// GET /health
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	components := h.checkAllComponents(c.Context())

	overallStatus := StatusUp
	for _, comp := range components {
		if comp.Status == StatusDown {
			overallStatus = StatusDown
			break
		}
	}

	response := HealthResponse{
		Status:     overallStatus,
		Version:    h.version,
		Uptime:     time.Since(h.startTime).Round(time.Second).String(),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Components: components,
	}

	statusCode := fiber.StatusOK
	if overallStatus == StatusDown {
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(response)
}

// Liveness indicates if the application is running
// GET /health/live
// Used by Kubernetes liveness probe
func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
	return c.JSON(LivenessResponse{
		Alive:     true,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Readiness indicates if the application is ready to receive traffic
// GET /health/ready
// Used by Kubernetes readiness probe
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
	components := h.checkAllComponents(c.Context())

	ready := true
	for _, comp := range components {
		if comp.Status == StatusDown {
			ready = false
			break
		}
	}

	response := ReadinessResponse{
		Ready:      ready,
		Components: components,
	}

	statusCode := fiber.StatusOK
	if !ready {
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(response)
}

// SystemInfo returns system information
// GET /health/system
func (h *HealthHandler) SystemInfo(c *fiber.Ctx) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	info := SystemInfo{
		GoVersion:     runtime.Version(),
		NumCPU:        runtime.NumCPU(),
		NumGoroutine:  runtime.NumGoroutine(),
		MemAlloc:      formatBytes(m.Alloc),
		MemTotalAlloc: formatBytes(m.TotalAlloc),
		MemSys:        formatBytes(m.Sys),
		NumGC:         m.NumGC,
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"system":  info,
		"uptime":  time.Since(h.startTime).Round(time.Second).String(),
		"version": h.version,
	})
}

// checkAllComponents checks all service dependencies concurrently
func (h *HealthHandler) checkAllComponents(ctx context.Context) map[string]ComponentHealth {
	components := make(map[string]ComponentHealth)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Check database
	wg.Add(1)
	go func() {
		defer wg.Done()
		health := h.checkDatabase(ctx)
		mu.Lock()
		components["database"] = health
		mu.Unlock()
	}()

	// Check Redis
	wg.Add(1)
	go func() {
		defer wg.Done()
		health := h.checkRedis(ctx)
		mu.Lock()
		components["redis"] = health
		mu.Unlock()
	}()

	wg.Wait()
	return components
}

// checkDatabase checks PostgreSQL connection
func (h *HealthHandler) checkDatabase(ctx context.Context) ComponentHealth {
	if h.db == nil {
		return ComponentHealth{
			Status:  StatusDown,
			Message: "database not configured",
		}
	}

	start := time.Now()

	sqlDB, err := h.db.DB()
	if err != nil {
		return ComponentHealth{
			Status:  StatusDown,
			Message: "failed to get database connection: " + err.Error(),
		}
	}

	// Use context with timeout for the ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return ComponentHealth{
			Status:  StatusDown,
			Message: "database ping failed: " + err.Error(),
		}
	}

	return ComponentHealth{
		Status:  StatusUp,
		Message: "postgresql connected",
		Latency: time.Since(start).Round(time.Microsecond).String(),
	}
}

// checkRedis checks Redis connection
func (h *HealthHandler) checkRedis(ctx context.Context) ComponentHealth {
	if h.redis == nil {
		return ComponentHealth{
			Status:  StatusDown,
			Message: "redis not configured",
		}
	}

	start := time.Now()

	// Use context with timeout for the ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := h.redis.Ping(pingCtx).Err(); err != nil {
		return ComponentHealth{
			Status:  StatusDown,
			Message: "redis ping failed: " + err.Error(),
		}
	}

	return ComponentHealth{
		Status:  StatusUp,
		Message: "redis connected",
		Latency: time.Since(start).Round(time.Microsecond).String(),
	}
}

// formatBytes formats bytes to human readable string
func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
