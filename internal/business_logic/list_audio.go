package business_logic

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

const limit = 5

func (bl *BusinessLogic) ListAudio(ctx context.Context, chatID int64, page int) (files []*models.File, hasPrevious bool, hasNext bool, err error) {
	audioFiles, hasNext, err := bl.repo.ListFilesByChatID(ctx, chatID, page, limit)

	if err != nil {
		return nil, false, false, fmt.Errorf("bl.repo.ListFilesByChatID :%w", err)
	}

	return audioFiles, page > 1, hasNext, nil
}
