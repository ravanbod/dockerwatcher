package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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
	m := make(map[string]interface{}) // I could not use normal fmt.Sprintf, because the message is multiline and mattermost returns error
	m["message"] = message
	m["channel_id"] = r.channelId
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+r.bearerAuth)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		err = errors.New("Status code is more than 399 error=" + strconv.Itoa(res.StatusCode) + " body=" + string(body))
	}
	return err
}
