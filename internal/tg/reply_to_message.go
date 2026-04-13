package tg

import (
	"context"
	"errors"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/logger"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) replyToMessage(ctx context.Context, chatID int64, messageForReplyId int, text string) (*tgBotApi.Message, error) {
	replyMessage := tgBotApi.NewMessage(chatID, text)
	replyMessage.ReplyToMessageID = messageForReplyId
	message, err := p.tgBotApi.Send(replyMessage)
	if err != nil {
		tgError, ok := errors.AsType[tgBotApi.Error](err)
		if ok {
			if tgError.Code == 403 {
				logger.FromContext(ctx).Info("bot was banned", "chatID", chatID)
			}
		}
		logger.FromContext(ctx).Error("replyToMessage", "error", err.Error())
		return nil, fmt.Errorf("p.tgBotApi.Send: %w", err)
	}
	return &message, nil
}
