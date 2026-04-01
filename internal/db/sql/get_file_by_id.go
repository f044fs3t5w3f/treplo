package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (r *repository) GetFileByID(ctx context.Context, fileID int64) (*models.File, error) {
	query := fmt.Sprintf(`
		SELECT %s 
		FROM files
		WHERE id = $1
	`, selectFields)
	row := r.db.QueryRowContext(ctx, query, fileID)

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("db.QueryRowContext: %w", err)
	}

	file := models.File{}
	if err := row.Scan(getFieldsForScan(&file)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("rows.Scan: %w", err)
	}
	return &file, nil
}
