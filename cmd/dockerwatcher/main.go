package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	dockerClient "github.com/docker/docker/client"
	"github.com/ravanbod/dockerwatcher/internal/config"
	"github.com/ravanbod/dockerwatcher/internal/repository/notification"
	"github.com/ravanbod/dockerwatcher/internal/repository/queue"
	"github.com/ravanbod/dockerwatcher/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appCfg, dockerCli, msgQueue := initApp(ctx)

	if appCfg.AppMode&config.WatcherApp != 0 {
		watcherService := service.NewWatcherService(dockerCli, msgQueue, appCfg.WatcherCfg.RedisQueueWriteName, appCfg.WatcherCfg.EventsFilter)
		slog.Info("Starting Watcher service ...")
		go watcherService.StartWatching(ctx)
	}

	if appCfg.AppMode&config.NotificationApp != 0 {
		var notifRepo notification.NotificationSender

		if appCfg.NotifCfg.NotificationPlatform == "generic" {
			slog.Info("Preparing Generic Notification Service")
			notifRepo, _ = notification.NewGenericNotificationSender(appCfg.NotifCfg.GenericNotifConfig.Url) // error always is nil
		} else if appCfg.NotifCfg.NotificationPlatform == "telegram" {
			slog.Info("Preparing Telegram Notification Service")
			notifRepo, _ = notification.NewTelegramNotificationSender(appCfg.NotifCfg.TelegramConfig.TelegramBotApiToken, appCfg.NotifCfg.TelegramConfig.TelegramChatID)
		} else if appCfg.NotifCfg.NotificationPlatform == "mattermost" {
			slog.Info("Preparing Mattermost Notification Service")
			notifRepo, _ = notification.NewMattermostNotificationSender(appCfg.NotifCfg.MattermostConfig.Host, appCfg.NotifCfg.MattermostConfig.BearerAuth, appCfg.NotifCfg.MattermostConfig.ChannelId)
		}
		notificationService := service.NewNotificationService(msgQueue, appCfg.NotifCfg.RedisQueueReadNames, notifRepo)
		slog.Info("Starting Notification service ...")
		go notificationService.StartListening(ctx)
	}

	<-ctx.Done()

	slog.Info("Shutting down in " + strconv.Itoa(int(appCfg.GracefulShutdownTimeout)) + " seconds")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(appCfg.GracefulShutdownTimeout)*time.Second)
	defer cancel()
	<-ctx.Done()
}

func initApp(ctx context.Context) (appCfg config.Config, dockerCli *dockerClient.Client, msgQueue queue.MessageQueue) {
	// Loading env vars
	appCfg, err := config.GetEnvConfig()
	if err != nil {
		slog.Error("Error in getting env config", "error", err)
		os.Exit(1)
	}

	if appCfg.AppMode == 0 {
		slog.Error("You have to enable either ENABLE_WATCHER or ENABLE_NOTIFICATION")
		os.Exit(1)
	}

	// Connect to the Docker engine
	if appCfg.AppMode&config.WatcherApp != 0 {
		cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithAPIVersionNegotiation())
		if err != nil {
			slog.Error("Error in connecting to the docker", "error", err)
			os.Exit(1)
		}
		_, err = cli.Ping(ctx)
		if err != nil {
			slog.Error("Error in connecting to the docker", "error", err)
			os.Exit(1)
		}
		dockerCli = cli
	}

	if appCfg.QueueConfig.QueueType == "redis" {
		// Connect to the Redis
		redisConn, err := queue.NewRedisClient(appCfg.QueueConfig.RedisConfig.RedisURL)

		if err != nil {
			slog.Error("Error in parsing the redis url", "error", err)
			os.Exit(1)
		}
		err = redisConn.Ping(ctx).Err()
		if err != nil {
			slog.Error("Error in connecting to the redis", "error", err)
			os.Exit(1)
		}
		msgQueue = queue.NewRedisRepo(redisConn)
	} else if appCfg.QueueConfig.QueueType == "dwqueue" {
		if appCfg.AppMode&config.WatcherApp == 0 || appCfg.AppMode&config.NotificationApp == 0 {
			slog.Error("For dwqueue, app mode must be notification and watcher")
			os.Exit(1)
		}
		msgQueue = queue.NewDWQueue(4096)
	}

	return
}
