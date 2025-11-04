package command

import (
	"context"
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository/sqlite"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/service"
)

type Handler interface {
	Process(ctx context.Context, msg *telegram.Message) (telegram.Chattable, error)
}

type Registry struct {
	tgBot      *telegram.BotAPI
	commandMap map[string]Handler
}

func NewRegistry(cfg *config.Config, tgBot *telegram.BotAPI) (*Registry, error) {
	db, err := sqlite.NewSQLiteDB(cfg.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("unable to create database instance: %w", err)
	}

	uow := sqlite.NewUnitOfWork(db)
	userService := service.NewUserService(cfg, uow)
	refuelService := service.NewRefuelService(cfg, uow)

	return &Registry{
		tgBot: tgBot,
		commandMap: map[string]Handler{
			"start": &startCommandHandler{
				users: userService,
			},
			"add": &addCommandHandler{
				users:   userService,
				refuels: refuelService,
			},
			// "stats": nil,
		},
	}, nil
}

func (r *Registry) ProcessCommand(ctx context.Context, msg *telegram.Message) (telegram.Chattable, error) {
	chatID := msg.Chat.ID

	handler, ok := r.commandMap[msg.Command()]
	if ok {
		return handler.Process(ctx, msg)
	} else {
		return createMessage(chatID, helpStartText), nil
	}
}

func createMessage(chatID int64, msgText string) telegram.Chattable {
	msg := telegram.NewMessage(chatID, telegram.EscapeText(telegram.ModeMarkdown, msgText))
	msg.ParseMode = telegram.ModeMarkdown

	return msg
}
