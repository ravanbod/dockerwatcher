package queue

import (
	"context"
	"errors"
	"log/slog"
)

// queueName is ignored because that app is not runned in multiple instances mode (multiple Watcher + 1 Notification)

type dwqueue struct {
	size  uint
	queue chan string
}

func NewDWQueue(size uint) MessageQueue {
	return dwqueue{size: size, queue: make(chan string, size)}
}

func (r dwqueue) PushMessageToQueue(ctx context.Context, queueName string, data string) error {
	r.queue <- data
	return nil
}

func (r dwqueue) GetLastMessageFromQueue(ctx context.Context, queueName string) (string, error) {
	select {
	case x, ok := <-r.queue:
		if ok {
			return x, nil
		} else {
			slog.Error("dwqueue channel is closed")
		}
	default:
	}
	return "", errors.New("QUEUE IS EMPTY")
}
