package bot

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
)

type App struct {
	botAPI          *telegram.BotAPI
	commandRegistry *CommandRegistry
}

func NewApp(cfg *config.Config) (*App, error) {
	httpClient := &http.Client{}

	if cfg.TelegramBot.ProxyAddress != "" {
		proxyURL, err := url.Parse(cfg.TelegramBot.ProxyAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
		}

		proxyDialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("failed to obtain proxy dialer: %w", err)
		}

		httpClient.Transport = &http.Transport{
			Dial: proxyDialer.Dial,
		}
	}

	botAPI, err := telegram.NewBotAPIWithClient(cfg.TelegramBot.Token, telegram.APIEndpoint, httpClient)
	if err != nil {
		return nil, fmt.Errorf("unable to create bot instance: %w", err)
	}
	botAPI.Debug = cfg.TelegramBot.Debug
	log.Printf("Authorized on account %s\n", botAPI.Self.UserName)

	commandRegistry, err := NewCommandRegistry(cfg, botAPI)
	if err != nil {
		return nil, fmt.Errorf("unable to create command registry: %w", err)
	}

	return &App{
		botAPI:          botAPI,
		commandRegistry: commandRegistry,
	}, nil
}

func (a *App) Run() {
	updates := a.botAPI.GetUpdatesChan(telegram.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 30,
	})

	for update := range updates {
		a.commandRegistry.ProcessUpdate(update)
	}
}
