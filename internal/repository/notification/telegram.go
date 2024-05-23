package notification

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramNotificationSender struct {
	botApi *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramNotificationSender(token string, charID int64) (NotificationSender, error) {
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return telegramNotificationSender{}, err
	}
	return telegramNotificationSender{botApi: botApi, chatID: charID}, nil
}

func (r telegramNotificationSender) SendMessage(message string) error {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered error from telegram", "error", r)
		}
	}()

	msg := tgbotapi.NewMessage(r.chatID, message)
	_, err := r.botApi.Send(msg)

	return err
}
