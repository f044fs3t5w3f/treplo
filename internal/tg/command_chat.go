package tg

import (
	"context"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) commandChat(ctx context.Context, update tgBotApi.Update, payload string) error {
	answer, err := p.businessLogic.AskAboutAudios(ctx, update.Message.Chat.ID, payload)
	var messageText string
	if err != nil {
		messageText = "Произошла ошибка"
	} else {
		messageText = answer
	}
	p.replyToMessage(
		ctx,
		update.Message.Chat.ID,
		update.Message.MessageID,
		messageText,
	)
	return nil
}
