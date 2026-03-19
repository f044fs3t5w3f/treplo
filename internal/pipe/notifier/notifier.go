package notifier

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

type Notifier interface {
	Notify(messageId int, chatId int64, message string) error
}

type NotifyProccessor struct {
	Notifier Notifier
}

func (n *NotifyProccessor) Process(ctx context.Context, file *models.File) error {
	if file.ProcessNotificationSent {
		return nil
	}

	err := n.Notifier.Notify(file.MessageID, file.ChatID, "Файл был обработан")
	if err != nil {
		return fmt.Errorf("n.Notifier.Notify: %w", err)
	}
	file.ProcessNotificationSent = true
	return nil
}
