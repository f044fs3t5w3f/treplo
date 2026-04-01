package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (r *repository) ListFilesByChatIDAndKeywords(ctx context.Context, keywords []string, chatID int64) ([]*models.File, error) {
	conditions, args := getConditionsAndArgs(keywords, chatID)
	query := fmt.Sprintf(`
		SELECT id, chat_id, message_id, file_id, filepath, salute_id,recognize_task_id, recognize_status, response_file_id, dialogue_content, process_notification_sent, encoding
		FROM files
		WHERE %s
		ORDER by id 
	`, conditions)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()
	files := make([]*models.File, 0)
	for rows.Next() {
		file := models.File{}
		if err := rows.Scan(&file.ID, &file.ChatID, &file.MessageID, &file.FileID, &file.Filepath, &file.SaluteId, &file.RecognizeTaskID, &file.RecognizeStatus, &file.ResponseFileID, &file.Content, &file.ProcessNotificationSent, &file.Encoding); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		files = append(files, &file)
	}
	return files, nil
}

func getConditionsAndArgs(keywords []string, chatID int64) (string, []any) {
	args := make([]any, len(keywords)+1)

	conditions := make([]string, len(keywords)+1)
	conditions[0] = "chat_id = $1"
	args[0] = chatID

	for i, keyword := range keywords {
		escapingKeyword := strings.ReplaceAll(keyword, `%`, `\%`)
		escapingKeyword = strings.ReplaceAll(escapingKeyword, `_`, `\_`)
		args[i+1] = "%" + escapingKeyword + "%"
		conditions[i+1] = fmt.Sprintf("dialogue_content ILIKE $%d", i+2)
	}
	return strings.Join(conditions, " AND "), args
}
