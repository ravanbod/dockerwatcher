package queue

import (
	"context"

	v9redis "github.com/redis/go-redis/v9"
)

func NewRedisClient(redisUrl string) (*v9redis.Client, error) {
	opts, err := v9redis.ParseURL(redisUrl)
	return v9redis.NewClient(opts), err
}

type redisQueue struct {
	rclient *v9redis.Client
}

func NewRedisRepo(rclient *v9redis.Client) MessageQueue {
	return redisQueue{rclient: rclient}
}

func (r redisQueue) PushMessageToQueue(ctx context.Context, queueName string, data string) error {
	return r.rclient.LPush(ctx, queueName, data).Err()
}

func (r redisQueue) GetLastMessageFromQueue(ctx context.Context, queueName string) (string, error) {
	data, err := r.rclient.RPop(ctx, queueName).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}
