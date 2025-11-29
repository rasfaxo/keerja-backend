package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	oauthStateKeyPrefix = "oauth:state:"
)

// ErrStateNotFound indicates the OAuth state token is missing or already consumed.
var ErrStateNotFound = errors.New("oauth state not found or expired")

// OAuthStateData captures metadata attached to an OAuth state token.
type OAuthStateData struct {
	RedirectURI          string    `json:"redirect_uri"`
	PostLoginRedirectURI string    `json:"post_login_redirect_uri,omitempty"`
	CodeChallenge        string    `json:"code_challenge,omitempty"`
	CodeChallengeMethod  string    `json:"code_challenge_method,omitempty"`
	ClientType           string    `json:"client_type,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

// OAuthStateStore abstracts state persistence so it can be backed by Redis or other stores.
type OAuthStateStore interface {
	Save(ctx context.Context, state string, data OAuthStateData, ttl time.Duration) error
	Consume(ctx context.Context, state string) (*OAuthStateData, error)
}

// RedisOAuthStateStore stores OAuth states in Redis with TTL.
type RedisOAuthStateStore struct {
	client *redis.Client
}

// NewRedisOAuthStateStore creates a new Redis-based OAuth state store.
func NewRedisOAuthStateStore(client *redis.Client) *RedisOAuthStateStore {
	return &RedisOAuthStateStore{client: client}
}

// Save implements OAuthStateStore.Save for Redis.
func (s *RedisOAuthStateStore) Save(ctx context.Context, state string, data OAuthStateData, ttl time.Duration) error {
	if s.client == nil {
		return errors.New("redis client is nil")
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal oauth state: %w", err)
	}

	if err := s.client.Set(ctx, oauthStateKeyPrefix+state, payload, ttl).Err(); err != nil {
		return fmt.Errorf("failed to store oauth state in redis: %w", err)
	}

	return nil
}

// Consume implements OAuthStateStore.Consume for Redis (get + delete).
func (s *RedisOAuthStateStore) Consume(ctx context.Context, state string) (*OAuthStateData, error) {
	if s.client == nil {
		return nil, errors.New("redis client is nil")
	}

	key := oauthStateKeyPrefix + state
	value, err := s.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrStateNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch oauth state from redis: %w", err)
	}

	if err := s.client.Del(ctx, key).Err(); err != nil {
		return nil, fmt.Errorf("failed to delete oauth state: %w", err)
	}

	var data OAuthStateData
	if err := json.Unmarshal(value, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal oauth state: %w", err)
	}

	return &data, nil
}

// InMemoryOAuthStateStore is primarily used in tests as a lightweight state store.
type InMemoryOAuthStateStore struct {
	mu     sync.RWMutex
	store  map[string]storedState
	ticker *time.Ticker
}

type storedState struct {
	Data      OAuthStateData
	ExpiresAt time.Time
}

// NewInMemoryOAuthStateStore creates a new in-memory store.
func NewInMemoryOAuthStateStore() *InMemoryOAuthStateStore {
	s := &InMemoryOAuthStateStore{
		store:  make(map[string]storedState),
		ticker: time.NewTicker(time.Minute),
	}

	go s.cleanupLoop()
	return s
}

// Save stores state in memory with TTL.
func (s *InMemoryOAuthStateStore) Save(_ context.Context, state string, data OAuthStateData, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[state] = storedState{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
	return nil
}

// Consume retrieves and deletes state entry.
func (s *InMemoryOAuthStateStore) Consume(_ context.Context, state string) (*OAuthStateData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.store[state]
	if !exists || time.Now().After(entry.ExpiresAt) {
		delete(s.store, state)
		return nil, ErrStateNotFound
	}

	delete(s.store, state)
	return &entry.Data, nil
}

func (s *InMemoryOAuthStateStore) cleanupLoop() {
	for range s.ticker.C {
		s.cleanup()
	}
}

func (s *InMemoryOAuthStateStore) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for state, entry := range s.store {
		if now.After(entry.ExpiresAt) {
			delete(s.store, state)
		}
	}
}
