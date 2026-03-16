package mechanics

import (
	"context"

	"github.com/a-kuleshov/treplo/internal/models"
)

type Repository interface {
	SaveFile(ctx context.Context, file *models.File) error
}
