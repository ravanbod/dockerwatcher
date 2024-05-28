package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
	res, err := http.Post(r.url, "application/json", bytes.NewBuffer(jsonByte))
	if res.StatusCode >= 400 {
		err = errors.New("Status code is more than 399 error=" + strconv.Itoa(res.StatusCode))
	}
	return err
}
