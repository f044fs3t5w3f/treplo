package sql

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (r *repository) ListFilesByChatID(ctx context.Context, chatID int64, page int, limit int) ([]*models.File, bool, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, chat_id, message_id, file_id, filepath, salute_id,recognize_task_id, recognize_status, response_file_id, dialogue_content, process_notification_sent, encoding, created_at
		FROM files
		WHERE chat_id = $1
		ORDER by id 
		LIMIT $2 OFFSET $3
	`, chatID, limit+1, (page-1)*limit)
	if err != nil {
		return nil, false, err
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}
	defer rows.Close()
	files := make([]*models.File, 0)
	for rows.Next() {
		file := models.File{}
		if err := rows.Scan(&file.ID, &file.ChatID, &file.MessageID, &file.FileID, &file.Filepath, &file.SaluteId, &file.RecognizeTaskID, &file.RecognizeStatus, &file.ResponseFileID, &file.Content, &file.ProcessNotificationSent, &file.Encoding, &file.CreatedAt); err != nil {
			return nil, false, fmt.Errorf("rows.Scan: %w", err)
		}
		files = append(files, &file)
	}
	hasNext := len(files) > limit
	if hasNext {
		files = files[:limit]
	}
	return files, hasNext, nil
}
