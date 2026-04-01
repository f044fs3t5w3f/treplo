package sql

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (r *repository) ListNewFiles(ctx context.Context) ([]*models.File, error) {
	query := fmt.Sprintf(`
		SELECT %s 
		FROM files
		where status = $1
	`, selectFields)
	rows, err := r.db.QueryContext(ctx, query, models.FileStatusNew)
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
