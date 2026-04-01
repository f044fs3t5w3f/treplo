package business_logic

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (bl *BusinessLogic) SaveUser(ctx context.Context, user *models.User) error {
	err := bl.repo.SaveUser(ctx, user)
	if err != nil {
		return fmt.Errorf("repo.SaveUser: %w", err)
	}
	return nil
}
