package redis

import (
	"context"

	v9redis "github.com/redis/go-redis/v9"
)

type WatcherRedisRepo struct {
	rclient   *v9redis.Client
	queueName string
}

func NewWatcherRedisRepo(rclient *v9redis.Client, queueName string) WatcherRedisRepo {
	return WatcherRedisRepo{rclient: rclient, queueName: queueName}
}

func (r *WatcherRedisRepo) PushMessageToQueue(ctx context.Context, data string) error {
	return r.rclient.LPush(ctx, r.queueName, data).Err()
}
