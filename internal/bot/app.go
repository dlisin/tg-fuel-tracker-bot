package bot

import (
	"fmt"
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
)

type App struct {
	botAPI          *telegram.BotAPI
	commandRegistry *CommandRegistry
}

func NewApp(cfg *config.Config) (*App, error) {
	botAPI, err := telegram.NewBotAPI(cfg.TelegramBot.Token)
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
