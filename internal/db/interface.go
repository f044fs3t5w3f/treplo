package db

import (
	"context"

	"github.com/a-kuleshov/treplo/internal/models"
)

type Repository interface {
	SaveFile(ctx context.Context, file *models.File) error
	ListFilesByChatID(ctx context.Context, chatID int64, page int, limit int) ([]*models.File, bool, error)
	ListFilesByChatIDAndKeywords(ctx context.Context, keywords []string, chatID int64) ([]*models.File, error)
	GetFileByID(ctx context.Context, fileID int64) (*models.File, error)
	ListNewFiles(ctx context.Context) ([]*models.File, error)
}
