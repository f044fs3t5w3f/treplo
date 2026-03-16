package mechanics

import (
	"context"
	"fmt"
)

type Mechanics struct {
	repo Repository
}

func NewMechanics(repo Repository) *Mechanics {
	return &Mechanics{repo: repo}
}

func (m *Mechanics) AddAudio(ctx context.Context, chatID int64, fileID string) error {
	err := m.repo.AddAudio(ctx, chatID, fileID)
	if err != nil {
		return fmt.Errorf("repo.AddAudio: %w", err)
	}
	return nil
}
