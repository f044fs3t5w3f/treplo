package tg

import (
	"context"
	"log/slog"
	"strings"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) onText(ctx context.Context, update tgBotApi.Update) error {
	message := update.Message
	command, payload := extractCommandAndPayload(message.Text)
	switch command {
	// list, get, find, chat
	case "start":
	// TODO: handle start command
	case "list":
		// TODO: handle list command
	case "get":
		// TODO: handle get command
	case "find":
		// TODO: handle find command
	case "chat":
		// TODO: handle chat command
	default:
		// TODO: handle unknown command
	}
	slog.DebugContext(ctx, "onText", "command", command, "payload", payload)

	return nil
}

func extractCommandAndPayload(text string) (string, string) {
	if !strings.HasPrefix(text, "/") {
		return "", ""
	}

	parts := strings.SplitN(text[1:], " ", 2)
	command := parts[0]
	var payload string
	if len(parts) > 1 {
		payload = strings.TrimSpace(parts[1])

	}
	return command, payload
}
