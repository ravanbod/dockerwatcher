package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockerClient "github.com/docker/docker/client"
	"github.com/ravanbod/dockerwatcher/internal/repository/queue"
)

type WatcherService struct {
	dockerCli    *dockerClient.Client
	msgQueue     queue.MessageQueue
	queueName    string
	eventFilters []string
}

func NewWatcherService(dockerCli *dockerClient.Client, msgQueue queue.MessageQueue, queueName string, eventFilters []string) WatcherService {
	return WatcherService{dockerCli: dockerCli, msgQueue: msgQueue, queueName: queueName, eventFilters: eventFilters}
}

// Blocking function! run it as a goroutine
func (r *WatcherService) StartWatching(ctx context.Context) {
	filterArgs := filters.NewArgs()
	for _, eventFilter := range r.eventFilters {
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
			err = r.msgQueue.PushMessageToQueue(ctx, r.queueName, string(jsonMsg))
			if err != nil {
				slog.Error("Error in pushing message to redis", "error", err)
			}
		case <-ctx.Done():
			slog.Info("Exiting Watcher service ...")
			return
		}
	}
}
