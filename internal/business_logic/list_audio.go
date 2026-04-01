package business_logic

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

const limit = 5

func (bl *BusinessLogic) ListAudio(ctx context.Context, chatID int64, page int) (files []*models.File, hasPrevious bool, hasNext bool, err error) {
	audioFiles, err := bl.repo.ListFilesByChatID(ctx, chatID)
	// TODO: put pagination into repo

	if err != nil {
		return nil, false, false, fmt.Errorf("bl.repo.ListFilesByChatID :%w", err)
	}
	totalCount := len(audioFiles)
	if len(audioFiles) <= limit {
		return audioFiles, false, false, nil
	}
	to := limit * page
	if to > totalCount {
		to = totalCount
	}
	audioFiles = audioFiles[(page-1)*limit : to]
	return audioFiles, page > 1, totalCount > limit*page, nil
}
