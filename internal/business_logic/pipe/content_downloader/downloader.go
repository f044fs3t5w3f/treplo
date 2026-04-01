package content_downloader

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/errors"
	"github.com/a-kuleshov/treplo/internal/models"
)

type Downloader interface {
	DownloadContent(ctx context.Context, saluteFileId string) (string, error)
}

type Tasker struct {
	Downloader Downloader
}

func (r *Tasker) Process(ctx context.Context, file *models.File) error {
	if file.Content != nil {
		return nil
	}
	if file.ResponseFileID == nil {
		return fmt.Errorf("%w: ResponseFileID", errors.ErrNoField)
	}
	content, err := r.Downloader.DownloadContent(ctx, *file.ResponseFileID)
	if err != nil {
		// TODO: check if error from Downloader is retriable and set file.Status = models.FileStatusError depending on it
		return fmt.Errorf("Downloader.DownloadContent: %w", err)
	}
	file.Content = &content
	file.Status = models.FileStatusDone
	return err
}
