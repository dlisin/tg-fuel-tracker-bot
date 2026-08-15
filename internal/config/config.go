package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

type StorageProvider string

const (
	StorageProviderSQLite StorageProvider = "sqlite"
)

type CacheProviderConfig string

const (
	CacheProviderMemory CacheProviderConfig = "memory"
)

type Config struct {
	Bot     BotConfig     `yaml:"bot"`
	Storage StorageConfig `yaml:"storage"`
	Cache   CacheConfig   `yaml:"cache"`
	Log     LogConfig     `yaml:"log"`
}

type BotConfig struct {
	Token string `yaml:"token" env:"TOKEN" required:"true"`

	DefaultCurrency    string   `yaml:"defaultCurrency" default:"₽"`
	DefaultFuelType    string   `yaml:"defaultFuelType" default:"ДТ"`
	AvailableFuelTypes []string `yaml:"availableFuelTypes" default:"АИ-92,АИ-95,АИ-100,ДТ"`
}

type StorageConfig struct {
	Provider StorageProvider     `yaml:"provider" env:"STORAGE_PROVIDER" default:"sqlite"`
	SQLite   SQLiteStorageConfig `yaml:"sqlite"`
}

type SQLiteStorageConfig struct {
	Path           string `yaml:"path" env:"STORAGE_PATH" default:"./fuelbot.db"`
	MaxConnections int    `yaml:"maxConnections" env:"STORAGE_MAX_CONNECTIONS" default:"1"`
}

type CacheConfig struct {
	Provider CacheProviderConfig `yaml:"provider" env:"CACHE_PROVIDER" default:"memory"`
	Memory   MemoryCacheConfig   `yaml:"memory"`
}

type MemoryCacheConfig struct {
	TTL time.Duration `yaml:"ttl" env:"CACHE_TTL" default:"30m"`
}

type LogConfig struct {
	Level  string `yaml:"level" env:"LOG_LEVEL" default:"info"`
	Format string `yaml:"format" env:"LOG_FORMAT" default:"text"`
}

func Load(logger *slog.Logger) (*Config, error) {
	logger = logger.With(slog.String("operation", "config.Load"))

	logger.Info("configuration loading started")

	var cfg Config

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		Files:     []string{"./local.env.yaml", "./config/env.yaml", "/etc/fuelbot/config.yaml"},
		EnvPrefix: "FUELBOT",
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".yml":  aconfigyaml.New(),
		},
	})

	if err := loader.Load(); err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate configuration: %w", err)
	}

	logger.Debug("configuration loaded", slog.Any("config", cfg))

	logger.Info("configuration loading completed")

	return &cfg, nil
}

func (cfg Config) Validate() error {
	return nil
}
