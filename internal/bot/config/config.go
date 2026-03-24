package config

import (
	"fmt"
	"log"

	commonConfig "github.com/kittipat1413/go-common/framework/config"
)

const (
	telegramBotToken = "TELEGRAM_BOT_TOKEN"
	telegramBotDebug = "TELEGRAM_BOT_DEBUG"
	databasePath     = "DATABASE_PATH"
	proxyAddress     = "PROXY_ADDRESS"

	defaultFuelType = "DEFAULT_FUEL_TYPE"
	defaultCurrency = "DEFAULT_CURRENCY"
)

type Config struct {
	TelegramBot TelegramBotConfig
	Database    DatabaseConfig

	DefaultFuelType string
	DefaultCurrency string
}

type TelegramBotConfig struct {
	Token string
	Debug bool

	ProxyAddress string
}

type DatabaseConfig struct {
	Path string
}

func Load() (*Config, error) {
	commonCfg := commonConfig.MustConfig(
		commonConfig.WithOptionalConfigPaths("./local.env.yaml", "./config/env.yaml", "/etc/fuelbot/config.yaml"),
		commonConfig.WithDefaults(map[string]any{
			telegramBotToken: "",
			telegramBotDebug: false,
			proxyAddress:     "",
			databasePath:     "/var/lib/fuelbot/fuelbot.db",
			defaultFuelType:  "ДТ",
			defaultCurrency:  "₽",
		}),
	)

	cfg := &Config{
		TelegramBot: TelegramBotConfig{
			Token:        commonCfg.GetString(telegramBotToken),
			Debug:        commonCfg.GetBool(telegramBotDebug),
			ProxyAddress: commonCfg.GetString(proxyAddress),
		},
		Database: DatabaseConfig{
			Path: commonCfg.GetString(databasePath),
		},
		DefaultFuelType: commonCfg.GetString(defaultFuelType),
		DefaultCurrency: commonCfg.GetString(defaultCurrency),
	}
	log.Printf("Loaded configuration: %+v\n", cfg)

	// Validate config
	if cfg.TelegramBot.Token == "" {
		return nil, fmt.Errorf("%s must be set", telegramBotToken)
	}

	return cfg, nil
}
