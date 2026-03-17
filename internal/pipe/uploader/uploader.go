package uploader

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/a-kuleshov/treplo/internal/models"
)

const directory = "storage"

type Uploader interface {
	UploadFile(file io.Reader) (string, error)
}

type FileUploader struct {
	Uploader Uploader
}

func (u *FileUploader) Process(ctx context.Context, file *models.File) error {
	saluteId, err := u.uploadFile(file)
	if err != nil {
		return fmt.Errorf("u.uploadFile: %w", err)
	}
	file.SaluteId = &saluteId
	return nil
}

func (u *FileUploader) uploadFile(file *models.File) (string, error) {
	fullFilename := filepath.Join(directory, *file.Filepath)
	f, err := os.Open(fullFilename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fileId, err := u.Uploader.UploadFile(f)
	if err != nil {
		return "", fmt.Errorf("uploader.UploadFile: %w", err)
	}
	return fileId, nil

}
