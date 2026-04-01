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
		SELECT %s
		FROM files
		WHERE %s
		ORDER by id 
	`, selectFields, conditions)

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

		if err := rows.Scan(getFieldsForScan(&file)...); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		files = append(files, &file)
	}
	return files, nil
}

func getConditionsAndArgs(keywords []string, chatID int64) (string, []any) {
	args := make([]any, len(keywords)+2)

	conditions := make([]string, len(keywords)+2)
	conditions[0] = "chat_id = $1 and status = $2"
	args[0] = chatID
	args[1] = models.FileStatusDone

	for i, keyword := range keywords {
		escapingKeyword := strings.ReplaceAll(keyword, `%`, `\%`)
		escapingKeyword = strings.ReplaceAll(escapingKeyword, `_`, `\_`)
		args[i+2] = "%" + escapingKeyword + "%"
		conditions[i+2] = fmt.Sprintf("dialogue_content ILIKE $%d", i+2)
	}
	return strings.Join(conditions, " AND "), args
}
