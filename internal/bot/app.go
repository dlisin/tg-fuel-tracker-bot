package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/command"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	tgBot           *telegram.BotAPI
	commandRegistry *command.Registry
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}
	log.Printf("Loaded configuration: %+v\n", cfg)

	tgBot, err := telegram.NewBotAPI(cfg.TelegramAPIToken)
	if err != nil {
		return nil, fmt.Errorf("unable to create bot instance: %w", err)
	}
	log.Printf("Authorized on account %s\n", tgBot.Self.UserName)
	tgBot.Debug = cfg.TelegramBotDebug

	commandRegistry, err := command.NewRegistry(cfg, tgBot)
	if err != nil {
		return nil, fmt.Errorf("unable to create command registry: %w", err)
	}

	return &App{
		tgBot:           tgBot,
		commandRegistry: commandRegistry,
	}, nil
}

func (a *App) Run() {
	updates := a.tgBot.GetUpdatesChan(telegram.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 30,
	})

	for update := range updates {
		a.handleUpdate(update)
	}
}

func (a *App) handleUpdate(update telegram.Update) {
	msg := update.Message
	ctx := context.Background()

	if msg != nil {
		log.Printf("Received message: %+v\n", msg)

		if msg.IsCommand() {
			replyMsg, err := a.commandRegistry.ProcessCommand(ctx, msg)
			if err != nil {
				log.Println("Failed to process message: ", err)
			}

			if replyMsg != nil {
				a.reply(replyMsg)
			}
		}
	}
}
func (a *App) reply(msg telegram.Chattable) {
	_, err := a.tgBot.Send(msg)
	if err != nil {
		log.Println("Failed to send message: ", err)
	}
}
