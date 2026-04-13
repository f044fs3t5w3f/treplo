package sql

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (r *repository) ListFilesByChatID(ctx context.Context, chatID int64, page int, limit int) ([]*models.File, bool, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM files
		WHERE chat_id = $1 and status = $2
		ORDER by id 
		LIMIT $3 OFFSET $4
	`, selectFields)
	rows, err := r.db.QueryContext(ctx, query, chatID, models.FileStatusDone, limit+1, (page-1)*limit)
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
		if err := rows.Scan(getFieldsForScan(&file)...); err != nil {
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
