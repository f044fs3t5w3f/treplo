package pipe

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
	"github.com/a-kuleshov/treplo/internal/pipe/downloader"
	"github.com/a-kuleshov/treplo/internal/pipe/recognizer"
	"github.com/a-kuleshov/treplo/internal/pipe/uploader"
	"github.com/a-kuleshov/treplo/pkg/sber/salute"
)

type repository interface {
	SaveFile(ctx context.Context, file *models.File) error
}

type FileProcessor interface {
	Process(ctx context.Context, file *models.File) error
}

type processors struct {
	processors []FileProcessor
	repo       repository
}

func NewPipe(repo repository, getFileURL downloader.GetFileURLfunc, saluteApi *salute.SpeachService) processors {
	return processors{
		processors: []FileProcessor{
			downloader.NewDownloader(getFileURL),
			&uploader.FileUploader{Uploader: saluteApi},
			&recognizer.Recognizer{Recognizer: saluteApi},
		},
		repo: repo,
	}
}

func (p processors) Process(ctx context.Context, file *models.File) error {
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
