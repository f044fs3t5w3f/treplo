package tg

import (
	"context"
	"strings"

	"github.com/a-kuleshov/treplo/internal/logger"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) onText(ctx context.Context, update tgBotApi.Update) error {
	message := update.Message
	command, payload := extractCommandAndPayload(message.Text)
	switch command {
	case "start":
		return p.commandStart(ctx, update)
	case "list":
		return p.commandList(ctx, update)
	case "find":
		return p.commandFind(ctx, update, payload)
	case "chat":
		return p.commandChat(ctx, update, payload)
	default:
		p.replyToMessage(
			ctx,
			update.Message.Chat.ID,
			update.Message.MessageID,
			"Неизвесная команда",
		)
	}
	logger.FromContext(ctx).Debug("onText", "command", command, "payload", payload)

	return nil
}

func extractCommandAndPayload(text string) (string, string) {
	if !strings.HasPrefix(text, "/") {
		return "", ""
	}

	parts := strings.SplitN(text[1:], " ", 2)
	command := strings.ToLower(parts[0])
	var payload string
	if len(parts) > 1 {
		payload = strings.TrimSpace(parts[1])

	}
	return command, payload
}

// Just send hello to reply
// Actualy we save all the information about user and chat for every interaction
func (p *Processor) commandStart(ctx context.Context, update tgBotApi.Update) error {
	p.replyToMessage(
		ctx,
		update.Message.Chat.ID,
		update.Message.MessageID,
		"Здравствуй, дорогой друг!",
	)
	return nil
}
