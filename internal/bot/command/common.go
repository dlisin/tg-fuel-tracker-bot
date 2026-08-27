package command

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/CloudyKit/jet/v6"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/template"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/service"
)

type commonCommand struct {
	botAPI  *telegram.Bot
	service service.BotService
}

func (h *commonCommand) sendMessageFromTemplate(ctx context.Context, chatID int64, templateName string, variables jet.VarMap) error {
	msgText, err := template.Render(templateName, variables)
	if err != nil {
		return err
	}

	return h.sendMessage(ctx, chatID, msgText)
}

func (h *commonCommand) sendMessage(ctx context.Context, chatID int64, msgText string) error {
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

func (h *commonCommand) resolveCar(ctx context.Context, userID domain.TelegramID, regNumber *domain.RegNumber) (*domain.Car, error) {
	cars, err := h.service.GetUserCars(ctx, userID)
	if err != nil {
		return nil, h.handleServiceError(err)
	}

	if regNumber != nil {
		for _, car := range cars {
			if car.RegNumber == *regNumber {
				return &car, nil
			}
		}

		return nil, fmt.Errorf("⚠️ Автомобиль с госномером %s не найден", *regNumber)
	}

	switch len(cars) {
	case 0:
		return nil, fmt.Errorf("⚠️ Сначала добавьте автомобиль")
	case 1:
		return &cars[0], nil
	default:
		return nil, fmt.Errorf("⚠️ У вас несколько автомобилей. Укажите госномер автомобиля")
	}
}

func (h *commonCommand) handleServiceError(err error) error {
	switch {
	case errors.Is(err, service.ErrCarNotFound):
		return errors.New("⚠️ Автомобиль не найден")

	case errors.Is(err, service.ErrCarAlreadyExists):
		return errors.New("⚠️ Автомобиль с таким госномером уже существует")

	case errors.Is(err, service.ErrRefuelNotFound):
		return errors.New("⚠️ Заправка не найдена")

	case errors.Is(err, service.ErrRefuelOdometerTooLow):
		return errors.New("⚠️ Пробег должен быть больше текущего пробега автомобиля")

	case errors.Is(err, service.ErrStatsNotEnoughRefuels):
		return errors.New("⚠️ Недостаточно заправок для расчета статистики")

	default:
		return errors.New("❌ Не удалось выполнить операцию. Попробуйте позже")
	}
}

func parseCommandArgs(text string) string {
	index := strings.IndexFunc(text, unicode.IsSpace)
	if index == -1 {
		return ""
	}

	return strings.TrimSpace(text[index:])
}
