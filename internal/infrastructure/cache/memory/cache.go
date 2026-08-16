package memory

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
)

type entry struct {
	value     any
	expiresAt time.Time
}

func (e entry) expired() bool {
	return time.Now().After(e.expiresAt)
}

type MemoryCache struct {
	logger  *slog.Logger
	mu      sync.Mutex
	entries map[string]entry
	ttl     time.Duration
}

func New(logger *slog.Logger, cfg config.MemoryCacheConfig) *MemoryCache {
	return &MemoryCache{
		logger: logger.With(
			slog.String("component", "MemoryCache"),
		),
		entries: make(map[string]entry),
		ttl:     cfg.TTL,
	}
}

func (c *MemoryCache) Get(ctx context.Context, key string) (any, error) {
	logger := c.logger.With(
		slog.String("operation", "Get"),
		slog.String("key", key),
	)

	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		logger.DebugContext(ctx, "entry missing")
		return nil, nil
	}

	if entry.expired() {
		delete(c.entries, key)

		logger.DebugContext(ctx, "entry expired")

		return nil, nil
	}
	logger.DebugContext(ctx, "entry found")

	return entry.value, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, value any) error {
	logger := c.logger.With(
		slog.String("operation", "Set"),
		slog.String("key", key),
		slog.Any("value", value),
	)

	c.mu.Lock()
	defer c.mu.Unlock()

	if value == nil {
		delete(c.entries, key)
		logger.DebugContext(ctx, "entry deleted")
	} else {
		c.entries[key] = entry{
			value:     value,
			expiresAt: time.Now().Add(c.ttl),
		}

		logger.DebugContext(ctx, "entry set")
	}

	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	logger := c.logger.With(
		slog.String("operation", "Delete"),
		slog.String("key", key),
	)

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
	logger.DebugContext(ctx, "entry deleted")

	return nil
}
