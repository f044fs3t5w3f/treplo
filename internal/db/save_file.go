package db

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (r *repository) SaveFile(ctx context.Context, file *models.File) error {
	if file.ID == 0 {
		result := r.db.QueryRowContext(ctx, `
			INSERT INTO 
			files (chat_id, file_id)
			VALUES ($1, $2)
			RETURNING id`, file.ChatID, file.FileID)
		err := result.Err()
		if err != nil {
			return fmt.Errorf("db.QueryRowContext: %w", err)
		}
		err = result.Scan(&file.ID)
		if err != nil {
			return fmt.Errorf("result.Scan: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			UPDATE files 
			SET 
				file_id = $1, 
				filepath = $2, 
				salute_id = $3,
				recognize_task_id = $4,
				recognize_status = $5,
				response_file_id = $6,
				dialogue_content = $7
			WHERE id = $8`,
			file.FileID, file.Filepath, file.SaluteId, file.RecognizeTaskID, file.RecognizeStatus, file.ResponseFileID, file.Content,
			file.ID)
		if err != nil {
			return fmt.Errorf("db.ExecContext: %w", err)
		}
	}
	return nil
}
