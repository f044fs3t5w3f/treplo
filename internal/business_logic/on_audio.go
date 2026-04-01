package business_logic

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (bl *BusinessLogic) OnAudio(ctx context.Context, chatID int64, messageId int, fileID string) error {
	// TODO: detect encoding
	file := models.File{
		FileID:    fileID,
		ChatID:    chatID,
		MessageID: messageId,
		Encoding:  nil,
	}
	err := bl.repo.SaveFile(ctx, &file)
	if err != nil {
		return fmt.Errorf("repo.SaveFile: %w", err)
	}

	bl.fileProcessor.Process(ctx, &file)

	return nil
}
