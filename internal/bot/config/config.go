package config

import (
	"fmt"

	commonConfig "github.com/kittipat1413/go-common/framework/config"
)

const (
	telegramBotToken = "TELEGRAM_BOT_TOKEN"
	telegramBotDebug = "TELEGRAM_BOT_DEBUG"
	databasePath     = "DATABASE_PATH"

	defaultFuelType = "DEFAULT_FUEL_TYPE"
	defaultCurrency = "DEFAULT_CURRENCY"
)

type Config struct {
	TelegramAPIToken string
	TelegramBotDebug bool

	DatabasePath string

	DefaultFuelType string
	DefaultCurrency string
}

func Load() (*Config, error) {
	cfg := commonConfig.MustConfig(
		commonConfig.WithOptionalConfigPaths("./local.env.yaml", "./config/env.yaml", "/etc/fuelbot/config.yaml"),
		commonConfig.WithDefaults(map[string]any{
			telegramBotToken: "",
			telegramBotDebug: false,
			databasePath:     "fuelbot.db",
			defaultFuelType:  "DT",
			defaultCurrency:  "₽",
		}),
	)

	if cfg.GetString(telegramBotToken) == "" {
		return nil, fmt.Errorf("%s must be set", telegramBotToken)
	}

	return &Config{
		TelegramAPIToken: cfg.GetString(telegramBotToken),
		TelegramBotDebug: cfg.GetBool(telegramBotDebug),
		DatabasePath:     cfg.GetString(databasePath),
		DefaultFuelType:  cfg.GetString(defaultFuelType),
		DefaultCurrency:  cfg.GetString(defaultCurrency),
	}, nil
}
