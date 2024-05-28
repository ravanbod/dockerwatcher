package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/ravanbod/dockerwatcher/internal/repository/notification"
	"github.com/ravanbod/dockerwatcher/internal/repository/redis"
	"github.com/ravanbod/dockerwatcher/pkg/jsontomd"
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
		qi := queueIndex % r.redisRepo.QueuesSize

		data, err := r.redisRepo.GetLastDataWithIndex(ctx, qi)
		if err == nil { // data available
			slog.Info("Reading nth queue", "n", qi, "data", data)
			messageText, err := jsontomd.ConvertJsonToMD(data)
			if err != nil {
				slog.Error("Failed to convert message to md...converting to tree", "error", err)
				messageText = jsontotree.ConvertJsonToTree(data)
			}

			err = r.notifRepo.SendMessage(messageText)
			if err != nil { // Error in sending message to notification platform ... resend the message to redis
				slog.Error("Error in sending message to notification platform", "error", err)
				slog.Info("Resending the message to redis", "message", data)
				r.redisRepo.PushMessageToQueue(ctx, qi, data)
			}
		}
		select {
		case <-time.After(time.Microsecond * 1000):
			queueIndex++
		case <-ctx.Done():
			slog.Info("Exiting Notification service ...")
			return
		}
	}
}
