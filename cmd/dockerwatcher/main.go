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
	"github.com/ravanbod/dockerwatcher/internal/repository/redis"
	"github.com/ravanbod/dockerwatcher/internal/service"
	v9redis "github.com/redis/go-redis/v9"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appCfg, dockerCli, redisConn := initApp(ctx)

	if appCfg.AppMode&config.WatcherApp == 1 {
		redisWatcherRepo := redis.NewWatcherRedisRepo(redisConn, appCfg.WatcherCfg.RedisQueueWriteName)
		watcherService := service.NewWatcherService(dockerCli, redisWatcherRepo)

		slog.Info("Starting Watcher service ...")
		go watcherService.StartWatching(ctx, appCfg.WatcherCfg.EventsFilter)
	}
	<-ctx.Done()

	slog.Info("Shutting down in " + strconv.Itoa(int(appCfg.GracefulShutdownTimeout)) + " seconds")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(appCfg.GracefulShutdownTimeout)*time.Second)
	defer cancel()
	<-ctx.Done()
}

func initApp(ctx context.Context) (appCfg config.Config, dockerCli *dockerClient.Client, redisConn *v9redis.Client) {
	// Loading env vars
	appCfg, err := config.GetEnvConfig()
	if err != nil {
		slog.Error("Error in getting env config")
		os.Exit(1)
	}

	if appCfg.AppMode == 0 {
		slog.Error("You have to enable either ENABLE_WATCHER or ENABLE_NOTIFICATION")
		os.Exit(1)
	}

	// Connect to the Docker engine
	if appCfg.AppMode&config.WatcherApp == 1 {
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

	// Connect to the Redis
	redisConn, err = redis.NewRedisClient(appCfg.RedisURL)

	if err != nil {
		slog.Error("Error in parsing the redis url", "error", err)
		os.Exit(1)
	}
	err = redisConn.Ping(ctx).Err()
	if err != nil {
		slog.Error("Error in connecting to the redis", "error", err)
		os.Exit(1)
	}

	return
}
