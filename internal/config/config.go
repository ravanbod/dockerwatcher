package config

import (
	"os"
	"strings"
)

type (
	Config struct {
		RedisURL   string
		WatcherCfg WatcherConfig
		NotifCfg   NotifConfig
	}

	WatcherConfig struct {
		RedisQueueWriteName string
	}

	NotifConfig struct {
		RedisQueueReadNames []string // comma seperated
	}
)

func GetEnvConfig() Config {
	return Config{
		RedisURL: getEnvWithDefault("REDIS_URL", "redis://localhost:6379/0?protocol=3"),
		WatcherCfg: WatcherConfig{
			RedisQueueWriteName: getEnvWithDefault("REDIS_QUEUE_WRITE_NAME", "dockerwatcher"),
		},
		NotifCfg: NotifConfig{
			RedisQueueReadNames: strings.Split(getEnvWithDefault("REDIS_QUEUE_READ_NAMES", "dockerwatcher,watcherdocker"), ","),
		},
	}
}

func getEnvWithDefault(key string, defaultKey string) string {
	if val := os.Getenv(key); val == "" {
		return defaultKey
	} else {
		return val
	}
}
