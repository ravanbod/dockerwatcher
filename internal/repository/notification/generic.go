package notification

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type genericNotificationSender struct {
	url string
}

func NewGenericNotificationSender(url string) (NotificationSender, error) {
	return genericNotificationSender{url: url}, nil
}

func (r genericNotificationSender) SendMessage(message string) error {
	jsonData := fmt.Sprintf("{\"message\":%s\"}", message)

	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", r.url, bytes.NewBuffer([]byte(jsonData)))
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
