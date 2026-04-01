package tg

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/a-kuleshov/treplo/internal/business_logic"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) commandFind(ctx context.Context, update tgBotApi.Update, payload string) error {
	var message string
	results, err := p.businessLogic.FindFiles(ctx, payload, update.Message.Chat.ID)

	if err == nil {
		textParts := make([]string, 0, len(results))
		for _, result := range results {
			textPart := fmt.Sprintf("%d", result.File.ID)
			textParts = append(textParts, textPart)
		}
		message = strings.Join(textParts, "\n")
	} else {
		if errors.Is(err, business_logic.ErrNoFiles) {
			message = "Не удалось найти записей с такими словами"
		} else {
			message = "Произошла ошибка, попробуйте позже"
		}
	}

	p.replyToMessage(
		ctx,
		update.Message.Chat.ID,
		update.Message.MessageID,
		message,
	)
	return nil
}
