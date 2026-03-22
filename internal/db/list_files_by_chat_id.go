package db

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (r *repository) ListFilesByChatID(ctx context.Context, chatID int64) ([]models.File, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, chat_id, message_id, file_id, filepath, salute_id,recognize_task_id, recognize_status, response_file_id, dialogue_content, process_notification_sent, encoding
		FROM files
		WHERE chat_id = $1
		ORDER by id 
	`, chatID)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()
	files := make([]models.File, 0)
	for rows.Next() {
		file := models.File{}
		if err := rows.Scan(&file.ID, &file.ChatID, &file.MessageID, &file.FileID, &file.Filepath, &file.SaluteId, &file.RecognizeTaskID, &file.RecognizeStatus, &file.ResponseFileID, &file.Content, &file.ProcessNotificationSent, &file.Encoding); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		files = append(files, file)
	}
	return files, nil
}
