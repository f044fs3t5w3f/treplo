package pipe

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
	"github.com/a-kuleshov/treplo/internal/pipe/content_downloader"
	"github.com/a-kuleshov/treplo/internal/pipe/downloader"
	"github.com/a-kuleshov/treplo/internal/pipe/tasker"
	"github.com/a-kuleshov/treplo/internal/pipe/uploader"
	"github.com/a-kuleshov/treplo/internal/pipe/waiter"
	"github.com/a-kuleshov/treplo/pkg/sber/salute"
)

type pipe struct {
	processors []FileProcessor
	repo       repository
}

func NewPipe(repo repository, getFileURL downloader.GetFileURLfunc, saluteApi *salute.SpeachService) pipe {
	return pipe{
		processors: []FileProcessor{
			downloader.NewDownloader(getFileURL),
			&uploader.FileUploader{Uploader: saluteApi},
			&tasker.Tasker{Tasker: saluteApi},
			&waiter.Waiter{StatusChecker: saluteApi},
			&content_downloader.Tasker{Downloader: saluteApi},
		},
		repo: repo,
	}
}

func (p pipe) Process(ctx context.Context, file *models.File) error {
	for _, processor := range p.processors {
		if err := processor.Process(ctx, file); err != nil {
			return err
		}
		if err := p.repo.SaveFile(ctx, file); err != nil {
			return fmt.Errorf("repo.SaveFile: %w", err)
		}
	}
	return nil
}
