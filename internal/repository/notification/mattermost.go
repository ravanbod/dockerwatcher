package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
	m := make(map[string]interface{}, 0)
	m["message"] = message
	m["channel_id"] = r.channelId
	jsonByte, err := json.Marshal(m)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	req.Header.Set("Authorization", "Bearer "+r.bearerAuth)
	res, err := client.Do(req)
	if res.StatusCode >= 400 {
		err = errors.New("Status code is more than 399 error=" + strconv.Itoa(res.StatusCode))
	}
	return err
}
