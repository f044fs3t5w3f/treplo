package tg

import (
	"context"
	"fmt"
	"strings"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Send the list of audio sent in CurrentChat
func (p *Processor) commandList(ctx context.Context, update tgBotApi.Update) error {
	files, err := p.businessLogic.ListAudio(ctx, update.Message.Chat.ID)
	if err != nil {
		p.replyToMessage(
			ctx,
			update.Message.Chat.ID,
			update.Message.MessageID,
			"Произошла ошибка, попробуйте позже",
		)
	}
	if len(files) == 0 {
		p.replyToMessage(
			ctx,
			update.Message.Chat.ID,
			update.Message.MessageID,
			"Пока ничего нет",
		)

	}
	textParts := make([]string, 0, len(files))
	for _, file := range files {
		textPart := fmt.Sprintf("%d", file.ID)
		textParts = append(textParts, textPart)
	}
	messageText := strings.Join(textParts, "\n")
	p.replyToMessage(
		ctx,
		update.Message.Chat.ID,
		update.Message.MessageID,
		messageText,
	)
	return nil
}
