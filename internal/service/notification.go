package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/ravanbod/dockerwatcher/internal/repository/notification"
	"github.com/ravanbod/dockerwatcher/internal/repository/redis"
	"github.com/ravanbod/dockerwatcher/pkg/jsontotree"
)

type NotificationService struct {
	redisRepo redis.NotificationRedisRepo
	notifRepo notification.NotificationSender
}

func NewNotificationService(redisRepo redis.NotificationRedisRepo, notifRepo notification.NotificationSender) NotificationService {
	return NotificationService{redisRepo: redisRepo, notifRepo: notifRepo}
}

func (r *NotificationService) StartListening(ctx context.Context) {
	queueIndex := uint(0)
	for {
		data, err := r.redisRepo.GetLastDataWithIndex(ctx, queueIndex%r.redisRepo.QueuesSize)
		if err == nil { // data available
			slog.Info("Reading nth queue", "n", queueIndex%r.redisRepo.QueuesSize, "data", data)
			r.notifRepo.SendMessage(jsontotree.ConvertJsonToTree(data))
		}
		select {
		case <-time.After(time.Microsecond * 100):
			queueIndex++
		case <-ctx.Done():
			slog.Info("Exiting Notification service ...")
			return
		}
	}
}
