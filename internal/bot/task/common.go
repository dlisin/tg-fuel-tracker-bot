package task

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/CloudyKit/jet/v6"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/template"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type commonTask struct {
	logger  *slog.Logger
	botAPI  *telegram.Bot
	service service.BotService
}

func (h *commonTask) sendMessageFromTemplate(ctx context.Context, chatID int64, templateName string, variables jet.VarMap) error {
	msgText, err := template.Render(templateName, variables)
	if err != nil {
		return err
	}

	return h.sendMessage(ctx, chatID, msgText)
}

func (h *commonTask) sendMessage(ctx context.Context, chatID int64, msgText string) error {
	_, err := h.botAPI.SendMessage(ctx, &telegram.SendMessageParams{
		ChatID:    chatID,
		Text:      msgText,
		ParseMode: models.ParseModeMarkdownV1,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
