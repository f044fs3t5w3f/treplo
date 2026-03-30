package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/logger"
	"github.com/a-kuleshov/treplo/internal/models"
)

func (r *repository) GetFileByID(ctx context.Context, fileID int64) (*models.File, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, chat_id, message_id, file_id, filepath, salute_id,recognize_task_id, recognize_status, response_file_id, dialogue_content, process_notification_sent, encoding, created_at 
		FROM files
		WHERE id = $1
	`, fileID)

	if err := row.Err(); err != nil {
		logger.FromContext(ctx).Error("GetFileByID", "error", err.Error())
		return nil, fmt.Errorf("db.QueryRowContext: %w", err)
	}

	file := models.File{}
	if err := row.Scan(&file.ID, &file.ChatID, &file.MessageID, &file.FileID, &file.Filepath, &file.SaluteId, &file.RecognizeTaskID, &file.RecognizeStatus, &file.ResponseFileID, &file.Content, &file.ProcessNotificationSent, &file.Encoding, &file.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.FromContext(ctx).Error("GetFileByID", "error", err.Error())
		return nil, fmt.Errorf("rows.Scan: %w", err)
	}
	return &file, nil
}
