package redis

import (
	v9redis "github.com/redis/go-redis/v9"
)

func NewRedisClient(redisUrl string) (*v9redis.Client, error) {
	opts, err := v9redis.ParseURL(redisUrl)
	return v9redis.NewClient(opts), err
}
