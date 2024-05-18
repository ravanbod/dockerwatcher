package redis

import (
	v9redis "github.com/redis/go-redis/v9"
)

func NewRedisClient(redisUrl string) *v9redis.Client {
	opts, err := v9redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}

	return v9redis.NewClient(opts)
}
