package tg

import (
	"context"
	"fmt"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Send the list of audio sent in CurrentChat
func (p *Processor) commandList(ctx context.Context, update tgBotApi.Update) error {
	files, _, hasNext, err := p.businessLogic.ListAudio(ctx, update.Message.Chat.ID, 1)
	if err != nil {
		p.replyToMessage(
			ctx,
			update.Message.Chat.ID,
			update.Message.MessageID,
			"Произошла ошибка, попробуйте позже",
		)
		return fmt.Errorf("businessLogic.ListAudio: %w", err)
	}
	if len(files) == 0 {
		p.replyToMessage(
			ctx,
			update.Message.Chat.ID,
			update.Message.MessageID,
			"Пока ничего нет",
		)
		return nil
	}
	keyboard := makeAudioFilesKeyboard(files, 1, false, hasNext)

	msg := tgBotApi.NewMessage(update.Message.Chat.ID, "Выберете беседу")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = keyboard
	_, err = p.tgBotApi.Send(msg)
	return err

}
