package tg

import (
	"context"
	"fmt"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Processor struct {
	service  mechanicService
	tgBotApi *tgBotApi.BotAPI
}

func NewProcessor(ctx context.Context, service mechanicService, tgBotApi *tgBotApi.BotAPI) *Processor {
	return &Processor{
		service:  service,
		tgBotApi: tgBotApi,
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
	err := p.service.AddAudio(ctx, message.Chat.ID, fileID)
	if err != nil {
		return fmt.Errorf("service.AddAudio: %w", err)
	}
	fileConfig := tgBotApi.FileConfig{FileID: fileID}
	p.tgBotApi.GetFile(fileConfig)
	return nil
}
