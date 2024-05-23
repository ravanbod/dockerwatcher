package service

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/ravanbod/dockerwatcher/internal/repository/redis"
)

type Watcher struct {
	ctx       context.Context
	dockerCli *dockerClient.Client
	redisRepo redis.WatcherRedisRepo
}

func NewWatcherService(ctx context.Context, dockerCli *dockerClient.Client, redisRepo redis.WatcherRedisRepo) Watcher {
	return Watcher{ctx: ctx, dockerCli: dockerCli, redisRepo: redisRepo}
}

// Blocking function! run it as a goroutine
func (r *Watcher) StartWatching() {
	msgs, errs := r.dockerCli.Events(r.ctx, types.EventsOptions{})
	for {
		select {
		case err := <-errs:
			slog.Error("error in reading docker evenets", err)
		case msg := <-msgs:
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				slog.Error("error in converting docker event to json", err)
				slog.Error("Ignoring the message ...")
				break
			}
			slog.Info("Docker event: " + string(jsonMsg))
			err = r.redisRepo.PushMessageToQueue(r.ctx, string(jsonMsg))
			if err != nil {
				slog.Error("Error in pushing message to redis", err)
			}
		case <-r.ctx.Done():
			slog.Info("Exiting Watcher service ...")
			return
		}
	}
}
