package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
)

func New(cfg config.LogConfig) (*slog.Logger, error) {
	var level slog.Level

	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", cfg.Level, err)
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler

	switch cfg.Format {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)

	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)

	default:
		return nil, fmt.Errorf("unsupported log format: %s", cfg.Format)
	}

	return slog.New(handler), nil
}
