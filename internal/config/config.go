package config

import (
	"errors"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		QueueConfig             QueueConfig
		WatcherCfg              WatcherConfig
		NotifCfg                NotifConfig
		AppMode                 int
		GracefulShutdownTimeout int64
	}

	QueueConfig struct {
		QueueType   string
		RedisConfig RedisConfig
	}

	RedisConfig struct {
		RedisURL string
	}

	WatcherConfig struct {
		RedisQueueWriteName string
		EventsFilter        []string // comma seperated
	}

	NotifConfig struct {
		RedisQueueReadNames  []string // comma seperated
		NotificationPlatform string
		TelegramConfig       TelegramConfig
		GenericNotifConfig   GenericNotifConfig
		MattermostConfig     MattermostConfig
	}

	TelegramConfig struct {
		TelegramBotApiToken string
		TelegramChatID      int64
	}

	GenericNotifConfig struct {
		Url string
	}

	MattermostConfig struct {
		Host       string
		BearerAuth string
		ChannelId  string
	}
)

const (
	WatcherApp = 1 << iota
	NotificationApp
)

func GetEnvConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file, trying to get envs from the shell", "error", err)
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
		return Config{}, err
	}

	genericNotifUrl := ""
	telegramChatID := 0
	mattermostHost := ""
	mattermostBearerAuth := ""
	mattermostChannelId := ""

	if getEnvWithDefault("NOTIFICATION_PLATFORM", "generic") == "generic" {
		genericNotifUrl = getEnvWithDefault("GENERIC_NOTIFICATION_URL", "http://localhost:80/webhook")
		_, err := url.ParseRequestURI(genericNotifUrl)
		if err != nil {
			slog.Error("Invalid URL", "error", err, "url", genericNotifUrl)
			return Config{}, err
		}
	} else if getEnvWithDefault("NOTIFICATION_PLATFORM", "generic") == "telegram" {
		telegramChatID, err = strconv.Atoi(getEnvWithDefault("TELEGRAM_CHAT_ID", "0"))
		if err != nil {
			slog.Error("Error in reading TELEGRAM_CHAT_ID", "error", err)
			return Config{}, err
		}
	} else if getEnvWithDefault("NOTIFICATION_PLATFORM", "generic") == "mattermost" {
		mattermostHost = getEnvWithDefault("MATTERMOST_HOST", "https://mattermost.local")
		_, err := url.ParseRequestURI(mattermostHost)
		if err != nil {
			slog.Error("Invalid URL", "error", err, "url", mattermostHost)
			return Config{}, err
		}
		mattermostBearerAuth = getEnvWithDefault("MATTERMOST_BEARER_AUTH", "xxxx")
		mattermostChannelId = getEnvWithDefault("MATTERMOST_CHANNEL_ID", "xxxx")
	} else {
		slog.Error("Error in reading NOTIFICATION_PLATFORM", "error", err)
		return Config{}, errors.New("NOTIFICATION_PLATFORM must be set")
	}

	queueType := getEnvWithDefault("QUEUE_TYPE", "redis")
	if queueType != "redis" && queueType != "dwqueue" {
		return Config{}, errors.New("QUEUE_TYPE must be either redis or dwqueue")
	}

	return Config{
		QueueConfig: QueueConfig{QueueType: queueType,
			RedisConfig: RedisConfig{
				RedisURL: getEnvWithDefault("REDIS_URL", "redis://localhost:6379/0?protocol=3")},
		},
		WatcherCfg: WatcherConfig{
			RedisQueueWriteName: getEnvWithDefault("REDIS_QUEUE_WRITE_NAME", "dockerwatcher"),
			EventsFilter:        strings.Split(getEnvWithDefault("EVENTS_FILTER", ""), ","),
		},
		NotifCfg: NotifConfig{
			RedisQueueReadNames:  strings.Split(getEnvWithDefault("REDIS_QUEUE_READ_NAMES", "dockerwatcher,watcherdocker"), ","),
			NotificationPlatform: getEnvWithDefault("NOTIFICATION_PLATFORM", "telegram"),
			TelegramConfig:       TelegramConfig{TelegramBotApiToken: getEnvWithDefault("TELEGRAM_BOT_API_TOKEN", "xxxx"), TelegramChatID: int64(telegramChatID)},
			GenericNotifConfig:   GenericNotifConfig{Url: genericNotifUrl},
			MattermostConfig:     MattermostConfig{Host: mattermostHost, BearerAuth: mattermostBearerAuth, ChannelId: mattermostChannelId},
		},
		AppMode:                 appMode,
		GracefulShutdownTimeout: int64(gracefulShutdownTimeout),
	}, nil
}

func getEnvWithDefault(key string, defaultKey string) string {
	if val := os.Getenv(key); val == "" {
		return defaultKey
	} else {
		return val
	}
}
