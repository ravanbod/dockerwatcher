package redis

import (
	"context"

	v9redis "github.com/redis/go-redis/v9"
)

type redisRepo struct {
	rclient *v9redis.Client
}

func NewRedisRepo(rclient *v9redis.Client) redisRepo {
	return redisRepo{rclient: rclient}
}

func (r *redisRepo) PushMessageToQueue(ctx context.Context, queueName string, data string) error {
	return r.rclient.LPush(ctx, queueName, data).Err()
}
