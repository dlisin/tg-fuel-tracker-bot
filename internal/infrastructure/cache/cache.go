package cache

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/infrastructure/cache/memory"
)

type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any) error
	Delete(ctx context.Context, key string) error
}

func New(logger *slog.Logger, cfg config.CacheConfig) (Cache, error) {
	switch cfg.Provider {
	case config.CacheProviderMemory:
		return memory.New(logger, cfg.Memory), nil

	default:
		return nil, fmt.Errorf("unsupported cache provider: %s", cfg.Provider)
	}
}
