package tg

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/business_logic"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Processor is a type that process update wherever it was obtained from (both with polling and webhooks)
type Processor struct {
	businessLogic *business_logic.BusinessLogic
	// tgBotApi *tgBotApi.BotAPI
}

func NewProcessor(ctx context.Context, service *business_logic.BusinessLogic, tgBotApi *tgBotApi.BotAPI) *Processor {
	return &Processor{
		businessLogic: service,
		// tgBotApi: tgBotApi,
	}
}

func (p *Processor) Process(ctx context.Context, update tgBotApi.Update) error {
	message := update.Message
	if message == nil {
		return nil
	}
	var fileID string
	if message.Audio != nil {
		fileID = message.Audio.FileID
	} else if message.Voice != nil {
		fileID = message.Voice.FileID
	} else {
		return nil
	}
	err := p.businessLogic.SaveFile(ctx, message.Chat.ID, message.MessageID, fileID)
	if err != nil {
		return fmt.Errorf("service.AddAudio: %w", err)
	}
	return nil
}
