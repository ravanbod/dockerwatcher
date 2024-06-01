package queue

import "context"

type MessageQueue interface {
	PushMessageToQueue(ctx context.Context, queueName string, data string) error
	GetLastMessageFromQueue(ctx context.Context, queueName string) (string, error)
}
