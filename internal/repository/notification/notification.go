package notification

type NotificationSender interface {
	SendMessage(string) error
}
