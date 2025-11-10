package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/command"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository/sqlite"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandRegistry struct {
	handlersMap map[string]command.Handler
}

func NewCommandRegistry(cfg *config.Config, botAPI *telegram.BotAPI) (*CommandRegistry, error) {
	db, err := sqlite.NewSQLiteDB(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("unable to create database instance: %w", err)
	}

	uow := sqlite.NewUnitOfWork(db)

	return &CommandRegistry{
		handlersMap: map[string]command.Handler{
			"start": command.NewStartHandler(cfg, botAPI, uow),
			"add":   command.NewAddHandler(cfg, botAPI, uow),
			// "stats": nil,
			// "import": nil,
			// "export": nil,
		},
	}, nil
}

func (r *CommandRegistry) ProcessUpdate(update telegram.Update) {
	msg := update.Message
	ctx := context.Background()

	if msg != nil {
		if msg.IsCommand() {
			log.Printf("Received command: %+v\n", msg)

			err := r.processCommand(ctx, msg)
			if err != nil {
				log.Printf("unable to process command: %+v\n", err)
			}
		}
	}
}

func (r *CommandRegistry) processCommand(ctx context.Context, msg *telegram.Message) error {
	handler, ok := r.handlersMap[msg.Command()]
	if !ok {
		return fmt.Errorf("unsupported command: %s", msg.Command())
	}

	return handler.Process(ctx, msg)
}
