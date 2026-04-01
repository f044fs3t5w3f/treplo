package tg

import (
	"context"
	"fmt"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) onVoice(ctx context.Context, update tgBotApi.Update) error {
	message := update.Message
	err := p.businessLogic.OnVoice(ctx, message.Chat.ID, message.MessageID, message.Voice.FileID)
	if err != nil {
		return fmt.Errorf("service.AddAudio: %w", err)
	}
	return nil
}
