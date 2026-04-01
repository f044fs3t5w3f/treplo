package memory

import (
	"context"
	"time"

	"github.com/a-kuleshov/treplo/internal/models"
)

func (m *MemoryRepository) SaveUser(ctx context.Context, user *models.User) error {
	for _, u := range m.users {
		if u.ID == user.ID {
			u.FirstName = user.FirstName
			u.LastName = user.LastName
			return nil
		}
	}
	user.CreatedAt = time.Now()
	m.users = append(m.users, user)
	return nil
}
