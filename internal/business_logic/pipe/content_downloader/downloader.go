package content_downloader

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

type Downloader interface {
	DownloadContent(saluteFileId string) (string, error)
}

type Tasker struct {
	Downloader Downloader
}

func (r *Tasker) Process(ctx context.Context, file *models.File) error {
	if file.Content != nil {
		return nil
	}
	if file.ResponseFileID == nil {
		return fmt.Errorf("file.ResponseFileId is nil")
	}
	content, err := r.Downloader.DownloadContent(*file.ResponseFileID)
	if err != nil {
		return fmt.Errorf("Downloader.DownloadContent: %w", err)
	}
	file.Content = &content
	return err
}
