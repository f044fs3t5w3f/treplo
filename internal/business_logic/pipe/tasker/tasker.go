package tasker

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/errors"

	"github.com/a-kuleshov/treplo/internal/models"
)

type SpeechTasker interface {
	CreateRecognizeTask(saluteFileId string, encoding string) (string, string, error)
}

type Tasker struct {
	Tasker SpeechTasker
}

func (r *Tasker) Process(ctx context.Context, file *models.File) error {
	if file.RecognizeTaskID != nil {
		return nil
	}

	if file.SaluteId == nil {
		return fmt.Errorf("%w: SaluteId", errors.ErrNoField)
	}

	if file.Encoding == nil {
		return fmt.Errorf("%w: Encoding", errors.ErrNoField)
	}
	taskId, status, err := r.Tasker.CreateRecognizeTask(*file.SaluteId, *file.Encoding)
	if err != nil {
		return fmt.Errorf("recognizer.CreateRecognizeTask: %w", err)
	}
	file.RecognizeTaskID = &taskId
	file.RecognizeStatus = &status
	return nil
}
