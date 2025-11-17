package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"time"
)

// Cache defines the interface for caching operations
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	DeletePattern(pattern string)
	Clear()
	Stats() CacheStats
}

// CacheStats holds cache statistics
type CacheStats struct {
	Hits        int64     `json:"hits"`
	Misses      int64     `json:"misses"`
	HitRatio    float64   `json:"hit_ratio"`
	Size        int       `json:"size"`
	MaxSize     int       `json:"max_size"`
	Evictions   int64     `json:"evictions"`
	LastCleanup time.Time `json:"last_cleanup"`
}

// cacheEntry represents a single cache entry with expiration
type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

// InMemoryCache implements Cache interface with in-memory storage
type InMemoryCache struct {
	data            sync.Map
	maxSize         int
	cleanupInterval time.Duration

	// Statistics
	hits        int64
	misses      int64
	evictions   int64
	lastCleanup time.Time
	mu          sync.RWMutex

	stopCleanup chan struct{}
}

// NewInMemoryCache creates a new in-memory cache with cleanup goroutine
func NewInMemoryCache(maxSize int, cleanupInterval time.Duration) *InMemoryCache {
	cache := &InMemoryCache{
		maxSize:         maxSize,
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan struct{}),
		lastCleanup:     time.Now(),
	}

	// Start cleanup goroutine
	go cache.startCleanup()

	return cache
}

// Get retrieves a value from the cache
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	value, ok := c.data.Load(key)
	if !ok {
		c.incrementMisses()
		return nil, false
	}

	entry := value.(cacheEntry)

	// Check if entry has expired
	if time.Now().After(entry.expiration) {
		c.data.Delete(key)
		c.incrementMisses()
		return nil, false
	}

	c.incrementHits()
	return entry.value, true
}

// Set stores a value in the cache with TTL
func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	entry := cacheEntry{
		value:      value,
		expiration: time.Now().Add(ttl),
	}

	// Check cache size before adding
	if c.getCurrentSize() >= c.maxSize {
		c.evictOldest()
	}

	c.data.Store(key, entry)
}

// Delete removes a key from the cache
func (c *InMemoryCache) Delete(key string) {
	c.data.Delete(key)
}

// DeletePattern removes all keys matching a pattern (simple prefix matching)
// Pattern examples:
//   - "company:*" matches all keys starting with "company:"
//   - "companies:list:*" matches all company list cache keys
func (c *InMemoryCache) DeletePattern(pattern string) {
	// Convert wildcard pattern to prefix
	prefix := pattern
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix = pattern[:len(pattern)-1]
	}

	// Iterate and delete matching keys
	c.data.Range(func(key, value interface{}) bool {
		keyStr := key.(string)
		if len(keyStr) >= len(prefix) && keyStr[:len(prefix)] == prefix {
			c.data.Delete(key)
		}
		return true
	})
}

// Clear removes all entries from the cache
func (c *InMemoryCache) Clear() {
	c.data.Range(func(key, value interface{}) bool {
		c.data.Delete(key)
		return true
	})

	c.mu.Lock()
	c.hits = 0
	c.misses = 0
	c.evictions = 0
	c.mu.Unlock()
}

// Stats returns cache statistics
func (c *InMemoryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := float64(c.hits + c.misses)
	hitRatio := 0.0
	if total > 0 {
		hitRatio = (float64(c.hits) / total) * 100
	}

	return CacheStats{
		Hits:        c.hits,
		Misses:      c.misses,
		HitRatio:    hitRatio,
		Size:        c.getCurrentSize(),
		MaxSize:     c.maxSize,
		Evictions:   c.evictions,
		LastCleanup: c.lastCleanup,
	}
}

// Stop stops the cleanup goroutine (call when shutting down)
func (c *InMemoryCache) Stop() {
	close(c.stopCleanup)
}

// Private methods

func (c *InMemoryCache) incrementHits() {
	c.mu.Lock()
	c.hits++
	c.mu.Unlock()
}

func (c *InMemoryCache) incrementMisses() {
	c.mu.Lock()
	c.misses++
	c.mu.Unlock()
}

func (c *InMemoryCache) incrementEvictions() {
	c.mu.Lock()
	c.evictions++
	c.mu.Unlock()
}

func (c *InMemoryCache) getCurrentSize() int {
	size := 0
	c.data.Range(func(key, value interface{}) bool {
		size++
		return true
	})
	return size
}

func (c *InMemoryCache) evictOldest() {
	// Simple eviction: remove first expired entry found
	// In a production system, implement LRU eviction
	var oldestKey interface{}
	var oldestTime time.Time = time.Now().Add(24 * time.Hour)

	c.data.Range(func(key, value interface{}) bool {
		entry := value.(cacheEntry)
		if entry.expiration.Before(oldestTime) {
			oldestTime = entry.expiration
			oldestKey = key
		}
		return true
	})

	if oldestKey != nil {
		c.data.Delete(oldestKey)
		c.incrementEvictions()
	}
}

func (c *InMemoryCache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

func (c *InMemoryCache) cleanup() {
	now := time.Now()

	c.data.Range(func(key, value interface{}) bool {
		entry := value.(cacheEntry)
		if now.After(entry.expiration) {
			c.data.Delete(key)
		}
		return true
	})

	c.mu.Lock()
	c.lastCleanup = now
	c.mu.Unlock()
}

// Helper functions

// GenerateFilterHash creates a consistent hash from filter parameters
func GenerateFilterHash(filters map[string]interface{}) string {
	if len(filters) == 0 {
		return "default"
	}

	// Sort keys for consistent hashing
	keys := make([]string, 0, len(filters))
	for k := range filters {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build filter string
	var filterStr string
	for _, k := range keys {
		filterStr += fmt.Sprintf("%s:%v:", k, filters[k])
	}

	// Hash the filter string
	hash := md5.Sum([]byte(filterStr))
	return hex.EncodeToString(hash[:8]) // Use first 8 bytes
}

// GenerateCacheKey creates a cache key from components
func GenerateCacheKey(components ...interface{}) string {
	key := ""
	for i, component := range components {
		if i > 0 {
			key += ":"
		}
		key += fmt.Sprintf("%v", component)
	}
	return key
}
