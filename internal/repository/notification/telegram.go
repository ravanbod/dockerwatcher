package notification

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type telegramNotificationSender struct {
	botToken string
	chatID   int64
}

func NewTelegramNotificationSender(token string, charID int64) (NotificationSender, error) {
	return telegramNotificationSender{botToken: token, chatID: charID}, nil
}

func (r telegramNotificationSender) SendMessage(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", r.botToken, "sendMessage")
	text := fmt.Sprintf("{\"chat_id\":%d, \"text\":\"%s\"}", r.chatID, message)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(text)))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		slog.Error("Error from telegram", "error", respBody)
		return errors.New("STATUS CODE IS NOT 200")
	}
	return err
}
