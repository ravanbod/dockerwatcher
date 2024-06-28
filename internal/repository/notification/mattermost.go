package notification

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type mattermostNotificationSender struct {
	host       string
	bearerAuth string
	channelId  string
}

func NewMattermostNotificationSender(host string, bearerAuth string, channelId string) (NotificationSender, error) {
	return mattermostNotificationSender{host: host, bearerAuth: bearerAuth, channelId: channelId}, nil
}

func (r mattermostNotificationSender) SendMessage(message string) error {
	url := r.host + "/api/v4/posts"
	jsonData := fmt.Sprintf("{\"message\":%s, \"channel_id\":\"%s\"}", message, r.channelId)

	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	req.Header.Set("Authorization", "Bearer "+r.bearerAuth)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		err = errors.New("Status code is more than 399 error=" + strconv.Itoa(res.StatusCode))
	}
	return nil
}
