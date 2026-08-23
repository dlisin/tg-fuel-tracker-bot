package task

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"log/slog"
	"text/template"

	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

type commonTask struct {
	logger  *slog.Logger
	cfg     config.BotPreferencesConfig
	botAPI  *telegram.Bot
	service service.BotService
}

func (h *commonTask) sendMessageFromTemplate(ctx context.Context, chatID int64, templateName string, data interface{}) error {
	t, err := template.ParseFS(templatesFS, templateName)
	if err != nil {
		return err
	}

	out := new(bytes.Buffer)
	err = t.Execute(out, data)
	if err != nil {
		return err
	}

	return h.sendMessage(ctx, chatID, out.String())
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
