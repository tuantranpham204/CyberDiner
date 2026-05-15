package store

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenDenylist tracks JWT identifiers (jti) that have been revoked before
// their natural expiration. Implementations should respect the supplied TTL
// so stale entries are dropped automatically.
type TokenDenylist interface {
	Add(ctx context.Context, jti string, ttl time.Duration) error
	Contains(ctx context.Context, jti string) (bool, error)
}

type memoryDenylist struct {
	mu      sync.RWMutex
	entries map[string]time.Time
	stop    chan struct{}
}

// NewMemoryDenylist returns an in-process denylist suitable for development
// and tests. A background goroutine sweeps expired entries every 5 minutes.
// Production deployments should swap this for a Redis-backed implementation.
func NewMemoryDenylist() TokenDenylist {
	d := &memoryDenylist{
		entries: make(map[string]time.Time),
		stop:    make(chan struct{}),
	}
	go d.gcLoop()
	return d
}

func (m *memoryDenylist) Add(_ context.Context, jti string, ttl time.Duration) error {
	if jti == "" || ttl <= 0 {
		return nil
	}
	m.mu.Lock()
	m.entries[jti] = time.Now().Add(ttl)
	m.mu.Unlock()
	return nil
}

func (m *memoryDenylist) Contains(_ context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, nil
	}
	m.mu.RLock()
	expiresAt, ok := m.entries[jti]
	m.mu.RUnlock()
	if !ok {
		return false, nil
	}
	if time.Now().After(expiresAt) {
		return false, nil
	}
	return true, nil
}

// ---------------------------------------------------------------------------
// Redis-backed denylist (used in production / when REDIS_ADDR is configured)
// ---------------------------------------------------------------------------

type redisDenylist struct {
	client    *redis.Client
	keyPrefix string
}

// NewRedisDenylist wires a denylist on top of a running Redis instance. Each
// jti is stored as a key with native TTL so cleanup happens automatically.
func NewRedisDenylist(client *redis.Client, keyPrefix string) TokenDenylist {
	if keyPrefix == "" {
		keyPrefix = "jwt:denylist:"
	}
	return &redisDenylist{client: client, keyPrefix: keyPrefix}
}

func (r *redisDenylist) Add(ctx context.Context, jti string, ttl time.Duration) error {
	if jti == "" || ttl <= 0 {
		return nil
	}
	return r.client.Set(ctx, r.keyPrefix+jti, "1", ttl).Err()
}

func (r *redisDenylist) Contains(ctx context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, nil
	}
	n, err := r.client.Exists(ctx, r.keyPrefix+jti).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// NewRedisClient parses a redis://user:pass@host:port/db DSN, opens the
// connection, and pings to verify reachability. Returns an error if the URL
// is invalid or Redis is unreachable so callers can decide whether to fall
// back to in-memory.
func NewRedisClient(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return client, nil
}

func (m *memoryDenylist) gcLoop() {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			now := time.Now()
			m.mu.Lock()
			for k, v := range m.entries {
				if now.After(v) {
					delete(m.entries, k)
				}
			}
			m.mu.Unlock()
		case <-m.stop:
			return
		}
	}
}
