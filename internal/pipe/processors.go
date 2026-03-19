package pipe

import (
	"context"

	"github.com/a-kuleshov/treplo/internal/models"
)

type repository interface {
	SaveFile(ctx context.Context, file *models.File) error
}

type FileProcessor interface {
	Process(ctx context.Context, file *models.File) error
}
