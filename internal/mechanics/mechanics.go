package mechanics

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
	"github.com/a-kuleshov/treplo/internal/pipe"
)

type Mechanics struct {
	repo          Repository
	fileProcessor pipe.FileProcessor
}

func NewMechanics(repo Repository, fileProcessor pipe.FileProcessor) *Mechanics {
	return &Mechanics{
		repo:          repo,
		fileProcessor: fileProcessor,
	}
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

	m.fileProcessor.Process(ctx, &file)

	return nil
}
