package main

import (
	"context"
	"log/slog"
	"os"

	dockerClient "github.com/docker/docker/client"
	"github.com/ravanbod/dockerwatcher/internal/config"
	"github.com/ravanbod/dockerwatcher/internal/repository/redis"
	"github.com/ravanbod/dockerwatcher/internal/service"
	v9redis "github.com/redis/go-redis/v9"
)

func main() {
	var ctx context.Context = context.Background()
	var appCfg, dockerCli, redisConn = initApp(ctx)

	if appCfg.AppMode&config.WatcherApp == 1 {
		var redisWatcherRepo = redis.NewWatcherRedisRepo(redisConn, appCfg.WatcherCfg.RedisQueueWriteName)
		var watcherService = service.NewWatcherService(dockerCli, redisWatcherRepo)

		slog.Info("Starting Watcher service ...")
		go watcherService.StartWatching(ctx)
	}
	for {
		select {}
	}
}

func initApp(ctx context.Context) (appCfg config.Config, dockerCli *dockerClient.Client, redisConn *v9redis.Client) {
	// Loading env vars
	appCfg = config.GetEnvConfig()
	if appCfg.AppMode == 0 {
		slog.Error("You have to enable either ENABLE_WATCHER or ENABLE_NOTIFICATION")
		os.Exit(1)
	}

	// Connect to the Docker engine
	if appCfg.AppMode&config.WatcherApp == 1 {
		cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithAPIVersionNegotiation())
		if err != nil {
			slog.Error("Error in connecting to the docker", err)
			os.Exit(1)
		}
		_, err = cli.Ping(ctx)
		if err != nil {
			slog.Error("Error in connecting to the docker", err)
			os.Exit(1)
		}
		dockerCli = cli
	}

	// Connect to the Redis
	redisConn, err := redis.NewRedisClient(appCfg.RedisURL)

	if err != nil {
		slog.Error("Error in parsing the redis url", err)
		os.Exit(1)
	}
	err = redisConn.Ping(ctx).Err()
	if err != nil {
		slog.Error("Error in connecting to the redis", err)
		os.Exit(1)
	}

	return
}
