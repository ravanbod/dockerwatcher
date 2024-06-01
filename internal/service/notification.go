package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/ravanbod/dockerwatcher/internal/repository/notification"
	"github.com/ravanbod/dockerwatcher/internal/repository/queue"
	"github.com/ravanbod/dockerwatcher/pkg/jsontomd"
	"github.com/ravanbod/dockerwatcher/pkg/jsontotree"
)

type NotificationService struct {
	msgQueue   queue.MessageQueue
	queueNames []string
	notifRepo  notification.NotificationSender
}

func NewNotificationService(msgQueue queue.MessageQueue, queueNames []string, notifRepo notification.NotificationSender) NotificationService {
	return NotificationService{msgQueue: msgQueue, queueNames: queueNames, notifRepo: notifRepo}
}

func (r *NotificationService) StartListening(ctx context.Context) {
	queueIndex := uint(0)
	for {
		qi := queueIndex % uint(len(r.queueNames))

		data, err := r.msgQueue.GetLastMessageFromQueue(ctx, r.queueNames[qi])
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
				slog.Info("Resending the message to the queue", "message", data)
				r.msgQueue.PushMessageToQueue(ctx, r.queueNames[qi], data)
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
