package tg

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/a-kuleshov/treplo/internal/business_logic"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) onCallbackQuery(ctx context.Context, update tgBotApi.Update) error {
	if strings.HasPrefix(update.CallbackQuery.Data, "p") {
		page, err := strconv.Atoi(update.CallbackQuery.Data[1:])
		if err != nil {
			p.replyToMessage(
				ctx,
				update.Message.Chat.ID,
				0,
				"Произошла ошибка",
			)
			return nil
		}
		return p.updateList(ctx, page, update.CallbackQuery.Message)
	} else {
		fileID, err := strconv.ParseInt(update.CallbackQuery.Data, 10, 64)
		if err != nil {
			p.replyToMessage(
				ctx,
				update.Message.Chat.ID,
				0,
				"Произошла ошибка",
			)
			return nil
		}
		return p.sendContent(ctx, fileID, *update.CallbackQuery.Message)

	}
}

func (p *Processor) updateList(ctx context.Context, page int, message *tgBotApi.Message) error {
	files, hasPrevious, hasNext, err := p.businessLogic.ListAudio(ctx, message.Chat.ID, page)
	if err != nil {
		p.replyToMessage(
			ctx,
			message.Chat.ID,
			message.MessageID,
			"Произошла ошибка, попробуйте позже",
		)
		return fmt.Errorf("businessLogic.ListAudio: %w", err)
	}
	keyboard := makeAudioFilesKeyboard(files, page, hasPrevious, hasNext)
	edit := tgBotApi.NewEditMessageReplyMarkup(message.Chat.ID, message.MessageID, keyboard)
	_, err = p.tgBotApi.Request(edit)
	return err
}

func (p *Processor) sendContent(ctx context.Context, fileID int64, message tgBotApi.Message) error {
	content, err := p.businessLogic.GetAudioContent(ctx, fileID, message.Chat.ID)
	var messageText string
	if err == nil {
		messageText = content
	} else {
		if errors.Is(err, business_logic.ErrNoAudio) {
			messageText = "Нет такого аудио"
		} else if errors.Is(err, business_logic.ErrNotReady) {
			messageText = "Обработка ещё не закончена"
		} else {
			messageText = "Произошла ошибка"
		}
	}
	p.replyToMessage(
		ctx,
		message.Chat.ID,
		message.MessageID,
		messageText,
	)
	return nil
}
