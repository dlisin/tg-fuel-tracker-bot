package command

import (
	"context"

	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type startCommand struct {
	commonCommand
}

func NewStartCommand(cfg config.BotConfig, botAPI *telegram.Bot, service service.BotService) Handler {
	return &startCommand{
		commonCommand: commonCommand{
			cfg:     cfg,
			botAPI:  botAPI,
			service: service,
		},
	}
}

func (h *startCommand) Process(ctx context.Context, msg *models.Message) error {
	_, err := h.botAPI.SetMyCommands(ctx, &telegram.SetMyCommandsParams{
		Scope: &models.BotCommandScopeChat{
			ChatID: msg.Chat.ID,
		},
		Commands: []models.BotCommand{
			{
				Command:     "start",
				Description: "помощь",
			},
			{
				Command:     "refuel_add",
				Description: "добавить заправку",
			},
			{
				Command:     "refuel_delete",
				Description: "удалить заправку",
			},
			{
				Command:     "refuel_list",
				Description: "показать заправки за указанный период",
			},
			{
				Command:     "refuel_stats",
				Description: "показать статистику за указанный период",
			},
		},
	})
	if err != nil {
		return err
	}

	return h.sendMessageFromTemplate(ctx, msg.Chat.ID, "templates/start.tmpl", nil)
}
