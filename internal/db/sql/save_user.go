package sql

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (r *repository) SaveUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, first_name, last_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET first_name = $2, last_name = $3
		RETURNING created_at
	`
	result := r.db.QueryRowContext(ctx, query, user.ID, user.FirstName, user.LastName)
	err := result.Err()
	if err != nil {
		return fmt.Errorf("db.QueryRowContext: %w", err)
	}
	err = result.Scan(&user.CreatedAt)
	if err != nil {
		return fmt.Errorf("result.Scan: %w", err)
	}
	return nil
}
