package sql

import "github.com/a-kuleshov/treplo/internal/models"

const selectFields = `id, chat_id, message_id, file_id, filepath, salute_id,recognize_task_id, recognize_status, response_file_id, dialogue_content, process_notification_sent, encoding, created_at, status `

func getFieldsForScan(file *models.File) []any {
	return []any{&file.ID, &file.ChatID, &file.MessageID, &file.FileID, &file.Filepath, &file.SaluteId, &file.RecognizeTaskID, &file.RecognizeStatus, &file.ResponseFileID, &file.Content, &file.ProcessNotificationSent, &file.Encoding, &file.CreatedAt, &file.Status}
}
