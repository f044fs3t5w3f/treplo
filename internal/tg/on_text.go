package tg

import (
	"context"
	"fmt"
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
		return p.commandStart(ctx, update)
	case "list":
		return p.commandList(ctx, update)
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

// Just send hello to reply
// Actualy we save all the information about user and chat for every interaction
func (p *Processor) commandStart(ctx context.Context, update tgBotApi.Update) error {
	p.replyToMessage(
		update.Message.Chat.ID,
		update.Message.MessageID,
		"Здравствуй, дорогой друг!",
	)
	return nil
}

// Send the list of audio sent in CurrentChat
func (p *Processor) commandList(ctx context.Context, update tgBotApi.Update) error {
	files, err := p.businessLogic.ListAudio(ctx, update.Message.Chat.ID)
	if err != nil {
		p.replyToMessage(
			update.Message.Chat.ID,
			update.Message.MessageID,
			"Произошла ошибка, попробуйте позже",
		)
	}
	if len(files) == 0 {
		p.replyToMessage(
			update.Message.Chat.ID,
			update.Message.MessageID,
			"Пока ничего нет",
		)

	}
	textParts := make([]string, 0, len(files))
	for _, file := range files {
		textPart := fmt.Sprintf("%d", file.ID)
		textParts = append(textParts, textPart)
	}
	messageText := strings.Join(textParts, "\n")
	p.replyToMessage(
		update.Message.Chat.ID,
		update.Message.MessageID,
		messageText,
	)
	return nil
}

func (p *Processor) replyToMessage(chatID int64, messageForReplyId int, text string) (*tgBotApi.Message, error) {
	replyMessage := tgBotApi.NewMessage(chatID, text)
	replyMessage.ReplyToMessageID = messageForReplyId
	message, err := p.tgBotApi.Send(replyMessage)
	if err != nil {
		// TODO: log message
		return nil, fmt.Errorf("p.tgBotApi.Send: %w", err)
	}
	return &message, nil

}
