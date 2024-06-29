package notification

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
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
	jsonData := fmt.Sprintf("{\"chat_id\":%d, \"text\":\"%s\"}", r.chatID, message)

	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		err = errors.New("Status code is more than 399 error=" + strconv.Itoa(res.StatusCode))
	}
	return err
}
