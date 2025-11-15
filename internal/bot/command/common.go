package command

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/repository"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

type Handler interface {
	Process(ctx context.Context, msg *telegram.Message) error
}

type commonCommand struct {
	cfg    *config.Config
	botAPI *telegram.BotAPI
	uow    repository.UnitOfWork
}

func (h *commonCommand) sendMessageFromTemplate(chatID int64, templateName string, data interface{}) error {
	t, err := template.ParseFS(templatesFS, templateName)
	if err != nil {
		return err
	}

	out := new(bytes.Buffer)
	err = t.Execute(out, data)
	if err != nil {
		return err
	}

	return h.sendMessage(chatID, out.String())
}

func (h *commonCommand) sendMessage(chatID int64, msgText string) error {
	msg := telegram.NewMessage(chatID, msgText)
	msg.ParseMode = telegram.ModeMarkdown

	_, err := h.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
