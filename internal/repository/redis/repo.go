package redis

import (
	"context"

	"github.com/pkg/errors"
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

type NotificationRedisRepo struct {
	rclient    *v9redis.Client
	queueNames []string
	QueuesSize uint
}

func NewNotificationRedisRepo(rclient *v9redis.Client, queueNames []string) NotificationRedisRepo {
	return NotificationRedisRepo{rclient: rclient, queueNames: queueNames, QueuesSize: uint(len(queueNames))}
}

func (r *NotificationRedisRepo) GetLastDataWithIndex(ctx context.Context, i uint) (string, error) {
	if i >= uint(len(r.queueNames)) {
		return "", errors.New("Out of index!!!")
	}

	data, err := r.rclient.RPop(ctx, r.queueNames[i]).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

func (r *NotificationRedisRepo) PushMessageToQueue(ctx context.Context, i uint, data string) error {
	return r.rclient.LPush(ctx, r.queueNames[i], data).Err()
}
