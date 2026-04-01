package tg

import (
	"context"
	"fmt"
	"time"

	"github.com/a-kuleshov/treplo/internal/business_logic"
	"github.com/a-kuleshov/treplo/internal/logger"
	"github.com/a-kuleshov/treplo/internal/models"
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
	err := p.process(ctx, update)
	if err != nil {
		logger.FromContext(ctx).Error("Error while processing update", "error", err.Error())
	}
	return err
}

func (p *Processor) process(ctx context.Context, update tgBotApi.Update) error {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if update.CallbackQuery != nil {
		return p.onCallbackQuery(ctx, update)
	}
	message := update.Message
	if message == nil {
		return nil
	}
	err := p.saveUser(ctx, message.From)
	if err != nil {
		logger.FromContext(ctx).Error("Error while saveUser", "error", err.Error())
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

func (p *Processor) saveUser(ctx context.Context, tgUser *tgBotApi.User) error {
	user := &models.User{
		ID:        tgUser.ID,
		FirstName: tgUser.FirstName,
		LastName:  tgUser.LastName,
	}
	err := p.businessLogic.SaveUser(ctx, user)
	if err != nil {
		return fmt.Errorf("service.SaveUser: %w", err)
	}
	return nil
}
