package tg

import (
	"context"
	"errors"
	"strconv"

	"github.com/a-kuleshov/treplo/internal/business_logic"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ErrIncorrectChatId = errors.New("incorrect chatID")

// Get conntent of file
func (p *Processor) commandGet(ctx context.Context, update tgBotApi.Update, payload string) error {
	fileID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		p.replyToMessage(
			ctx,
			update.Message.Chat.ID,
			update.Message.MessageID,
			"Некорректный ID",
		)
		return ErrIncorrectChatId
	}
	content, err := p.businessLogic.GetAudioContent(ctx, fileID, update.Message.Chat.ID)
	var messageText string
	if err == nil {
		messageText = content
	} else {
		if errors.Is(err, business_logic.ErrNoAudio) {
			messageText = "Нет такого аудио"
		} else if errors.Is(err, business_logic.ErrNotReady) {
			messageText = "Обработка ещё не закончена"
		} else {
			messageText = "Произошла ошибка"
		}
	}
	p.replyToMessage(
		ctx,
		update.Message.Chat.ID,
		update.Message.MessageID,
		messageText,
	)
	return nil
}
