package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-kuleshov/treplo/internal/models"
)

type GetFileURLfunc func(fileID string) (string, error)

type Downloader struct {
	getFileURL  GetFileURLfunc
	storagePath string
}

func NewDownloader(getFileURL GetFileURLfunc, storagePath string) (*Downloader, error) {
	// TODO: create directory if not exists
	return &Downloader{
		getFileURL:  getFileURL,
		storagePath: storagePath,
	}, nil
}

func (d *Downloader) Process(ctx context.Context, file *models.File) error {
	if file.Filepath != nil {
		return nil
	}
	return d.Download(ctx, file)
}

func (d *Downloader) Download(ctx context.Context, file *models.File) error {
	if file.Filepath != nil {
		return nil
	}
	url, err := d.getFileURL(file.FileID)
	if err != nil {
		return err
	}

	reader, err := downloadFile(url)
	if err != nil {
		return err
	}
	defer reader.Close()

	filename := fmt.Sprintf("%d_%s", file.ID, file.FileID)
	fullFilename := filepath.Join(d.storagePath, filename)

	f, err := os.OpenFile(fullFilename, os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}
	file.Filepath = &filename
	return nil
}

func downloadFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http.Get: %s", resp.Status)
	}
	return resp.Body, nil
}
