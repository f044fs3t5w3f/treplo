package uploader

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/errors"
	"github.com/a-kuleshov/treplo/internal/models"
)

type Uploader interface {
	UploadFile(ctx context.Context, file io.Reader) (string, error)
}

type FileUploader struct {
	Uploader    Uploader
	StoragePath string
}

func (u *FileUploader) Process(ctx context.Context, file *models.File) error {
	if file.SaluteId != nil {
		return nil
	}
	saluteId, err := u.uploadFile(ctx, file)
	if err != nil {
		return err
	}
	file.SaluteId = &saluteId
	return nil
}

func (u *FileUploader) uploadFile(ctx context.Context, file *models.File) (string, error) {
	if file.Filepath == nil {
		return "", fmt.Errorf("%w: Filepath", errors.ErrNoField)
	}
	fullFilename := filepath.Join(u.StoragePath, *file.Filepath)
	f, err := os.Open(fullFilename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fileId, err := u.Uploader.UploadFile(ctx, f)
	if err != nil {
		return "", fmt.Errorf("uploader.UploadFile: %w", err)
	}
	return fileId, nil

}
