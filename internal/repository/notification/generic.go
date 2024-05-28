package notification

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type genericNotificationSender struct {
	url string
}

func NewGenericNotificationSender(url string) (NotificationSender, error) {
	return genericNotificationSender{url: url}, nil
}

func (r genericNotificationSender) SendMessage(message string) error {
	m := make(map[string]interface{}, 0)
	m["message"] = message
	jsonByte, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = http.Post(r.url, "application/json", bytes.NewBuffer(jsonByte))
	return err
}
