package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		RedisURL                string
		WatcherCfg              WatcherConfig
		NotifCfg                NotifConfig
		AppMode                 int
		GracefulShutdownTimeout int64
	}

	WatcherConfig struct {
		RedisQueueWriteName string
		EventsFilter        []string // comma seperated
	}

	NotifConfig struct {
		RedisQueueReadNames []string // comma seperated
	}
)

const (
	WatcherApp = 1 << iota
	NotificationApp
)

func GetEnvConfig() Config {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}

	appMode := 0
	if getEnvWithDefault("ENABLE_WATCHER", "0") == "1" || strings.ToLower(getEnvWithDefault("ENABLE_WATCHER", "0")) == "true" {
		appMode = appMode | WatcherApp
	}

	if getEnvWithDefault("ENABLE_NOTIFICATION", "0") == "1" || strings.ToLower(getEnvWithDefault("ENABLE_NOTIFICATION", "0")) == "true" {
		appMode = appMode | NotificationApp
	}

	gracefulShutdownTimeout, err := strconv.Atoi(getEnvWithDefault("GRACEFUL_SHUTDOWN_TIMEOUT", "10"))
	if err != nil {
		slog.Error("Error in reading GRACEFUL_SHUTDOWN_TIMEOUT", "error", err)
	}

	return Config{
		RedisURL: getEnvWithDefault("REDIS_URL", "redis://localhost:6379/0?protocol=3"),
		WatcherCfg: WatcherConfig{
			RedisQueueWriteName: getEnvWithDefault("REDIS_QUEUE_WRITE_NAME", "dockerwatcher"),
			EventsFilter:        strings.Split(getEnvWithDefault("EVENTS_FILTER", ""), ","),
		},
		NotifCfg: NotifConfig{
			RedisQueueReadNames: strings.Split(getEnvWithDefault("REDIS_QUEUE_READ_NAMES", "dockerwatcher,watcherdocker"), ","),
		},
		AppMode:                 appMode,
		GracefulShutdownTimeout: int64(gracefulShutdownTimeout),
	}
}

func getEnvWithDefault(key string, defaultKey string) string {
	if val := os.Getenv(key); val == "" {
		return defaultKey
	} else {
		return val
	}
}
