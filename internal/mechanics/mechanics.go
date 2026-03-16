package mechanics

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

type Mechanics struct {
	repo Repository
}

func NewMechanics(repo Repository) *Mechanics {
	return &Mechanics{repo: repo}
}

func (m *Mechanics) SaveFile(ctx context.Context, chatID int64, fileID string) error {
	file := models.File{
		FileID: fileID,
		ChatID: chatID,
	}
	err := m.repo.SaveFile(ctx, &file)
	if err != nil {
		return fmt.Errorf("repo.SaveFile: %w", err)
	}
	return nil
}
