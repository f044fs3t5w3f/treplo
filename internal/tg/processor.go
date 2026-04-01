package tg

import (
	"context"
	"time"

	"github.com/a-kuleshov/treplo/internal/business_logic"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Processor is a type that process update wherever it was obtained from (both with polling and webhooks)
type Processor struct {
	businessLogic *business_logic.BusinessLogic
	tgBotApi      *tgBotApi.BotAPI
}

func NewProcessor(ctx context.Context, service *business_logic.BusinessLogic, tgBotApi *tgBotApi.BotAPI) *Processor {
	return &Processor{
		businessLogic: service,
		tgBotApi:      tgBotApi,
	}
}

func (p *Processor) Process(ctx context.Context, update tgBotApi.Update) error {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	p.saveUser(ctx, update)
	if update.CallbackQuery != nil {
		return p.onCallbackQuery(ctx, update)
	}
	message := update.Message
	if message == nil {
		return nil
	}
	if message.Voice != nil {
		return p.onVoice(ctx, update)
	} else if message.Audio != nil {
		return p.onAudio(ctx, update)
	} else if message.Text != "" {
		return p.onText(ctx, update)
	}
	return nil
}

func (p *Processor) saveUser(ctx context.Context, update tgBotApi.Update) error {
	// TODO: save user
	return nil
}
