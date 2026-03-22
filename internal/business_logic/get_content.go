package business_logic

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrNoAudio  = errors.New("no such audio")
	ErrNotReady = errors.New("audio was not processed")
)

func (bl *BusinessLogic) GetAudioContent(ctx context.Context, fileID int64, chatID int64) (string, error) {
	file, err := bl.repo.GetFileByID(ctx, fileID)
	if err != nil {
		return "", fmt.Errorf("bl.repo.GetFileByID: %w", err)
	}

	if file == nil || file.ChatID != chatID {
		return "", ErrNoAudio
	}
	if file.RecognizeStatus == nil || *file.RecognizeStatus != "DONE" || file.Content == nil {
		return "", ErrNotReady
	}

	return *file.Content, nil
}
