package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitRedis initializes a Redis client using configuration values.
func InitRedis(cfg *Config) (*redis.Client, error) {
	if cfg == nil {
		cfg = LoadConfig()
	}

	var opts *redis.Options
	var err error

	if cfg.RedisURL != "" {
		opts, err = redis.ParseURL(cfg.RedisURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse REDIS_URL: %w", err)
		}
	} else {
		opts = &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		}
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	redisClient = client
	return client, nil
}

// GetRedis returns the initialized Redis client.
func GetRedis() *redis.Client {
	if redisClient == nil {
		panic("Redis client not initialized. Call InitRedis() first")
	}
	return redisClient
}

// CloseRedis closes the Redis client connection gracefully.
func CloseRedis() error {
	if redisClient == nil {
		return nil
	}

	if err := redisClient.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	redisClient = nil
	return nil
}
