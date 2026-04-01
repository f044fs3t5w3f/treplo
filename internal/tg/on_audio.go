package tg

import (
	"context"
	"fmt"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) onAudio(ctx context.Context, update tgBotApi.Update) error {
	message := update.Message
	fileID := message.Audio.FileID
	err := p.businessLogic.OnAudio(ctx, message.Chat.ID, message.MessageID, fileID)
	if err != nil {
		return fmt.Errorf("service.AddAudio: %w", err)
	}
	return nil
}
