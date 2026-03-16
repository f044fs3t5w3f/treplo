package db

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (r *repository) SaveFile(ctx context.Context, file *models.File) error {
	result := r.db.QueryRowContext(ctx, `
INSERT INTO 
files (chat_id, file_id)
VALUES ($1, $2)
RRETURNING id`, file.ChatID, file.FileID)
	err := result.Err()
	if err != nil {
		return fmt.Errorf("db.QueryRowContext: %w", err)
	}
	err = result.Scan(&file.ID)
	return err
}
