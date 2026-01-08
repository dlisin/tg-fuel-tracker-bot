package command

import (
	"context"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/model"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

type addCommand struct {
	commonCommand
}

func NewAddCommand(cfg *config.Config, botAPI *telegram.BotAPI, refuelRepository repository.RefuelRepository) Handler {
	return &addCommand{
		commonCommand: commonCommand{
			cfg:              cfg,
			botAPI:           botAPI,
			refuelRepository: refuelRepository,
		},
	}
}

func (h *addCommand) Process(ctx context.Context, msg *telegram.Message) error {
	userID := model.TelegramID(msg.From.ID)
	prevRefuel, err := h.refuelRepository.GetByOdometer(ctx, userID, 0)
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "❌ Не удалось сохранить заправку. Попробуйте позже")
	}

	cmdArgs, err := parseAddCommandArgs(msg.CommandArguments(), prevRefuel)
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "⚠️ Ошибка ввода: "+err.Error())
	}

	newRefuel, err := h.refuelRepository.Create(ctx, &model.Refuel{
		UserID:        userID,
		Odometer:      cmdArgs.Odometer,
		Liters:        cmdArgs.Liters,
		PriceTotal:    cmdArgs.TotalPrice,
		PricePerLiter: cmdArgs.TotalPrice / cmdArgs.Liters,
		CreatedAt:     time.Now(),
	})
	if err != nil {
		return h.sendMessage(msg.Chat.ID, "❌ Не удалось сохранить заправку. Попробуйте позже")
	}

	var stats *model.RefuelStats
	if prevRefuel != nil {
		stats = model.CreateRefuelStats([]model.Refuel{*prevRefuel, *newRefuel})
	}

	return h.sendMessageFromTemplate(msg.Chat.ID, "templates/add.tmpl", struct {
		Refuel *model.Refuel
		Stats  *model.RefuelStats
		Config *config.Config
	}{
		Refuel: newRefuel,
		Stats:  stats,
		Config: h.cfg,
	})
}
