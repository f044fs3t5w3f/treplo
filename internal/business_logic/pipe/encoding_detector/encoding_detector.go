package encoding_detector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	pipeErrors "github.com/a-kuleshov/treplo/internal/business_logic/pipe/errors"
	"github.com/a-kuleshov/treplo/internal/logger"
	"github.com/a-kuleshov/treplo/internal/models"
)

const encodingMP3 = "MP3"

var ErrUnsupportedEncoding = pipeErrors.NewErrorForUser("unsupported encoding", "Неподдерживаемая кодировка файла")

type EncodingDetector struct {
	StoragePath string
}

func (ed EncodingDetector) Process(ctx context.Context, file *models.File) error {
	if file.Filepath == nil {
		return fmt.Errorf("%w: Filepath", pipeErrors.ErrNoField)
	}
	fullFilename := filepath.Join(ed.StoragePath, *file.Filepath)
	fileToRecogize, err := os.Open(fullFilename)
	if err != nil {
		if os.IsNotExist(err) {
			logger.FromContext(ctx).Warn("file doesn't exists", "filePath", *file.Filepath)
		}
		return fmt.Errorf("os.Open: %w", err)
	}
	defer fileToRecogize.Close()
	encoding, err := detectEncoding(fileToRecogize)
	if err != nil {
		return fmt.Errorf("detectEncoding: %w", err)
	}
	file.Encoding = &encoding
	return nil
}

func detectEncoding(file *os.File) (string, error) {
	bytes := make([]byte, 3)
	bytesRead, err := file.Read(bytes)
	if err != nil {
		// even in case of io.EOF
		return "", err
	}
	if bytesRead != len(bytes) {
		return "", ErrUnsupportedEncoding
	}
	if string(bytes) == "ID3" {
		return encodingMP3, nil
	}
	return "", ErrUnsupportedEncoding

}
