package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockerClient "github.com/docker/docker/client"
	"github.com/ravanbod/dockerwatcher/internal/repository/redis"
)

type WatcherService struct {
	dockerCli *dockerClient.Client
	redisRepo redis.WatcherRedisRepo
}

func NewWatcherService(dockerCli *dockerClient.Client, redisRepo redis.WatcherRedisRepo) WatcherService {
	return WatcherService{dockerCli: dockerCli, redisRepo: redisRepo}
}

// Blocking function! run it as a goroutine
func (r *WatcherService) StartWatching(ctx context.Context, eventFilters []string) {
	filterArgs := filters.NewArgs()
	for _, eventFilter := range eventFilters {
		if eventFilter == "" {
			break
		}
		slog.Info("Added filter", "filter", eventFilter)
		splited := strings.Split(eventFilter, "=")
		filterArgs.Add(splited[0], splited[1])
	}

	msgs, errs := r.dockerCli.Events(ctx, types.EventsOptions{Filters: filterArgs})
	for {
		select {
		case err := <-errs:
			slog.Error("error in reading docker evenets", "error", err)
		case msg := <-msgs:
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				slog.Error("error in converting docker event to json", "error", err)
				slog.Error("Ignoring the message ...")
				break
			}
			slog.Info("Docker event", "event", string(jsonMsg))
			err = r.redisRepo.PushMessageToQueue(ctx, string(jsonMsg))
			if err != nil {
				slog.Error("Error in pushing message to redis", "error", err)
			}
		case <-ctx.Done():
			slog.Info("Exiting Watcher service ...")
			return
		}
	}
}
