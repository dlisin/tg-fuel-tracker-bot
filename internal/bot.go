package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/db"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/stats"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App2 struct {
	Bot *tgbot.BotAPI
	DB  *db.DB
}

func (a *bot.App) Reply(chatID int64, text string) {
	msg := tgbot.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"

	_, _ = a.Bot.Send(msg)
}

func (a *bot.App) HandleUpdate(update tgbot.Update) {
	// callbacks
	if update.CallbackQuery != nil {
		cq := update.CallbackQuery
		fromID := cq.From.ID
		data := strings.TrimSpace(cq.Data)
		if strings.HasPrefix(data, "stats:") {
			arg := strings.TrimPrefix(data, "stats:")
			_ = a.handleStats(fromID, cq.Message.Chat.ID, arg)
			_, _ = a.Bot.Request(tgbot.NewCallback(cq.ID, ""))
		}
		return
	}

	if update.Message == nil {
		return
	}
	msg := update.Message
	chatID := msg.Chat.ID

	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			a.Reply(chatID, HelpText)
		case "register":
			params := strings.TrimSpace(msg.CommandArguments())
			carMake, fuelType, odo, err := parseRegister(params)
			if err != nil {
				a.Reply(chatID, "❌ Ошибка /register: "+err.Error())
				return
			}
			if err := a.DB.UpsertUser(msg.From.ID, carMake, fuelType, odo); err != nil {
				a.Reply(chatID, "❌ Не удалось сохранить профиль: "+err.Error())
				return
			}
			a.Reply(chatID, fmt.Sprintf("✅ Сохранено. %s, топливо: %s, одометр: %d", carMake, fuelType, odo))
		case "add":
			args := strings.TrimSpace(msg.CommandArguments())
			odo, liters, tp, ppl, err := parseAdd(args)
			if err != nil {
				a.Reply(chatID, "❌ Ошибка /add: "+err.Error())
				return
			}
			u, err := a.DB.GetUserByTG(msg.From.ID)
			if err != nil {
				a.Reply(chatID, "⚠️ Сначала зарегистрируйтесь: /register <марка>; <топливо>; <пробег>")
				return
			}
			lastOdo, err := a.DB.GetLastOdometer(u.ID)
			if err != nil {
				a.Reply(chatID, "❌ Ошибка проверки одометра: "+err.Error())
				return
			}
			if lastOdo > 0 && odo < lastOdo {
				a.Reply(chatID, fmt.Sprintf("❌ Одометр меньше предыдущего (%d). Проверьте ввод.", lastOdo))
				return
			}
			if err := a.DB.AddFillup(u.ID, odo, liters, tp, ppl); err != nil {
				a.Reply(chatID, "❌ Не удалось добавить заправку: "+err.Error())
				return
			}
			var priceInfo string
			if tp != nil {
				priceInfo = fmt.Sprintf("чек: %.2f", *tp)
			} else {
				priceInfo = fmt.Sprintf("цена/л: %.3f", *ppl)
			}
			a.Reply(chatID, fmt.Sprintf("⛽ Добавлено: одометр %d, %.2f л, %s", odo, liters, priceInfo))
		case "stats":
			arg := strings.TrimSpace(msg.CommandArguments())
			_ = a.handleStats(msg.From.ID, chatID, arg)
		default:
			a.Reply(chatID, "Неизвестная команда. Напишите /start")
		}
	} else {
		a.Reply(chatID, "Напишите /start для помощи, /register — для регистрации, /add — для добавления заправки.")
	}
}

func (a *bot.App) handleStats(fromUserID int64, chatID int64, arg string) error {
	u, err := a.DB.GetUserByTG(fromUserID)
	if err != nil {
		a.Reply(chatID, "⚠️ Сначала зарегистрируйтесь: /register <марка>; <топливо>; <пробег>")
		return err
	}
	rng, label, err := stats.ParseStatsRange(strings.TrimSpace(arg), time.Now())
	if err != nil {
		a.Reply(chatID, "❌ "+err.Error())
		return err
	}
	var startEnd *[2]time.Time
	if rng != nil {
		startEnd = rng
	}
	fills, err := a.DB.GetFillups(u.ID, startEnd)
	if err != nil {
		a.Reply(chatID, "❌ Ошибка выборки данных: "+err.Error())
		return err
	}
	st := stats.Calc(fills)
	if st.Entries < 2 || st.DistanceKm <= 0 {
		a.Reply(chatID, "ℹ️ Недостаточно данных. Нужны минимум две записи в выбранном периоде.")
		return nil
	}
	a.Reply(chatID, fmt.Sprintf(
		"📊 Статистика %s\n• Пробег: %.0f км\n• Средний расход: %.2f л/100км\n• Цена/л: %.3f → %.3f (%+.3f; %+.1f%%)",
		label, st.DistanceKm, st.AvgConsumption, st.FirstPPL, st.LastPPL, st.PriceDeltaAbs, st.PriceDeltaPct))

	return nil
}

// --- Parsers (kept here for simplicity) -----------------------------------------

func parseRegister(s string) (carMake, fuelType string, odo int64, err error) {
	if s == "" {
		return "", "", 0, fmt.Errorf("укажите: <марка>; <тип_топлива>; <пробег>")
	}
	parts := splitBySemicolon(s)
	if len(parts) < 3 {
		return "", "", 0, fmt.Errorf("нужно 3 параметра, разделённых точкой с запятой")
	}
	carMake = strings.TrimSpace(parts[0])
	fuelType = strings.TrimSpace(parts[1])
	odo64, err := parseInt64(strings.TrimSpace(parts[2]))
	if err != nil {
		return "", "", 0, fmt.Errorf("пробег должен быть числом")
	}
	if odo64 < 0 {
		return "", "", 0, fmt.Errorf("пробег должен быть ≥ 0")
	}
	return carMake, fuelType, odo64, nil
}

func parseAdd(s string) (odometer int64, liters float64, totalPrice *float64, pricePerLiter *float64, err error) {
	if s == "" {
		return 0, 0, nil, nil, fmt.Errorf("укажите: <пробег> <литры> <сумма_чека|цена_за_литр>")
	}
	fs := strings.Fields(s)
	if len(fs) < 3 {
		return 0, 0, nil, nil, fmt.Errorf("ожидалось 3 параметра")
	}
	odo, err := parseInt64(fs[0])
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("пробег должен быть целым числом")
	}
	if odo < 0 {
		return 0, 0, nil, nil, fmt.Errorf("пробег должен быть ≥ 0")
	}
	lit, err := parseFloat(fs[1])
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("литры должны быть числом")
	}
	third := fs[2]
	if strings.Contains(strings.ToLower(third), "/") || strings.HasSuffix(strings.ToLower(third), "l") {
		third = strings.TrimSuffix(strings.ToLower(third), "/l")
		ppl, err := parseFloat(third)
		if err != nil {
			return 0, 0, nil, nil, fmt.Errorf("цена/л должна быть числом")
		}
		return odo, lit, nil, &ppl, nil
	}
	tp, err := parseFloat(third)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("сумма чека должна быть числом")
	}
	return odo, lit, &tp, nil, nil
}

// tiny helpers
func splitBySemicolon(s string) []string {
	parts := strings.Split(s, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}
func parseInt64(s string) (int64, error) {
	s = strings.ReplaceAll(s, ",", "")
	return strconv.ParseInt(s, 10, 64)
}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.TrimSuffix(s, "l")
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}
