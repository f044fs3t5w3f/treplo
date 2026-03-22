package memory

import (
	"context"
	"strings"

	"github.com/a-kuleshov/treplo/internal/models"
)

// ListFilesByChatIDAndKeywords implements [db.Repository].
func (m *MemoryRepository) ListFilesByChatIDAndKeywords(ctx context.Context, keywords []string, chatID int64) ([]*models.File, error) {
	files := make([]*models.File, 0)
	for _, file := range m.files {
		if file.ChatID != chatID {
			continue
		}
		content := strings.ToLower(*file.Content)
		hasEntry := len(keywords) == 0
		for _, keyword := range keywords {
			if strings.Contains(content, strings.ToLower(keyword)) {
				hasEntry = true
				break
			}
		}
		if hasEntry {
			files = append(files, deepCopyFile(file))
		}
	}
	return files, nil
}
