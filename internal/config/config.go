package config

import "os"

type (
	Config struct {
		RedisURL       string
		RedisQueueName string
	}
)

func GetEnvConfig() Config {
	return Config{
		RedisURL:       getEnvWithDefault("REDIS_URL", "redis://localhost:6379/0?protocol=3"),
		RedisQueueName: getEnvWithDefault("REDIS_QUEUE_NAME", "dockerwatcher"),
	}
}

func getEnvWithDefault(key string, defaultKey string) string {
	if val := os.Getenv(key); val == "" {
		return defaultKey
	} else {
		return val
	}
}
