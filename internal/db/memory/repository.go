package memory

import (
	"context"
	"sync/atomic"

	"github.com/a-kuleshov/treplo/internal/db"
	"github.com/a-kuleshov/treplo/internal/models"
)

// MemoryRepository implements [db.Repository]. Only for testing purposes.
type MemoryRepository struct {
	files     []*models.File
	currentID atomic.Int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		files:     []*models.File{},
		currentID: atomic.Int64{},
	}
}

// GetFileByID implements [db.Repository].
func (m *MemoryRepository) GetFileByID(ctx context.Context, fileID int64) (*models.File, error) {
	for _, f := range m.files {
		if f.ID == fileID {
			return deepCopyFile(f), nil
		}
	}
	return nil, nil
}

// ListFiles implements [db.Repository].
func (m *MemoryRepository) ListNewFiles(ctx context.Context) ([]*models.File, error) {
	files := make([]*models.File, 0)
	for _, f := range m.files {
		if f.Status != models.FileStatusNew {
			continue
		}
		files = append(files, deepCopyFile(f))
	}
	return files, nil
}

// ListFilesByChatID implements [db.Repository].
func (m *MemoryRepository) ListFilesByChatID(ctx context.Context, chatID int64, page int, limit int) ([]*models.File, bool, error) {
	files := make([]*models.File, 0, limit)

	for _, f := range m.files {
		if f.ChatID == chatID {
			if len(files) == limit {
				return files, true, nil
			}
			files = append(files, deepCopyFile(f))
		}
	}
	return files, false, nil
}

// SaveFile implements [db.Repository].
func (m *MemoryRepository) SaveFile(ctx context.Context, file *models.File) error {
	if file.ID == 0 {
		m.currentID.Add(1)
		file.ID = m.currentID.Load()
		m.files = append(m.files, file)
	} else {
		for i, f := range m.files {
			if f.ID == file.ID {
				m.files[i] = file
				return nil
			}
		}
	}
	return nil
}

var _ db.Repository = &MemoryRepository{}

func strPtr(s *string) *string {
	if s == nil {
		return nil
	}
	v := *s
	return &v
}

func deepCopyFile(file *models.File) *models.File {
	if file == nil {
		return nil
	}

	return &models.File{
		ID:                      file.ID,
		FileID:                  file.FileID,
		ChatID:                  file.ChatID,
		MessageID:               file.MessageID,
		Filepath:                strPtr(file.Filepath),
		Encoding:                strPtr(file.Encoding),
		SaluteId:                strPtr(file.SaluteId),
		RecognizeTaskID:         strPtr(file.RecognizeTaskID),
		RecognizeStatus:         strPtr(file.RecognizeStatus),
		ResponseFileID:          strPtr(file.ResponseFileID),
		Content:                 strPtr(file.Content),
		ProcessNotificationSent: file.ProcessNotificationSent,
	}
}
