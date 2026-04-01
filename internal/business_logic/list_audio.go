package business_logic

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (bl *BusinessLogic) ListAudio(ctx context.Context, chatID int64) ([]*models.File, error) {
	audioFiles, err := bl.repo.ListFilesByChatID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("bl.repo.ListFilesByChatID :%w", err)
	}
	return audioFiles, nil
}
