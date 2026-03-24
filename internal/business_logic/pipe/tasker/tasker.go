package tasker

import (
	"context"
	"errors"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

type SpeechTasker interface {
	CreateRecognizeTask(saluteFileId string, encoding string) (string, string, error)
}

type Tasker struct {
	Tasker SpeechTasker
}

var ErrNoField = errors.New("required fields is missing")

func (r *Tasker) Process(ctx context.Context, file *models.File) error {
	if file.RecognizeTaskID != nil {
		return nil
	}

	if file.SaluteId == nil {
		return fmt.Errorf("%w: SaluteId", ErrNoField)
	}

	if file.Encoding == nil {
		return fmt.Errorf("%w: Encoding", ErrNoField)
	}
	taskId, status, err := r.Tasker.CreateRecognizeTask(*file.SaluteId, *file.Encoding)
	if err != nil {
		return fmt.Errorf("recognizer.CreateRecognizeTask: %w", err)
	}
	file.RecognizeTaskID = &taskId
	file.RecognizeStatus = &status
	return nil
}
