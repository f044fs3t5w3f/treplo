package tasker

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

type SpeechTasker interface {
	CreateRecognizeTask(saluteFileId string) (string, string, error)
}

type Tasker struct {
	Tasker SpeechTasker
}

func (r *Tasker) Process(ctx context.Context, file *models.File) error {
	return r.CreateRecognizeTask(ctx, file)
}

func (r *Tasker) CreateRecognizeTask(ctx context.Context, file *models.File) error {
	taskId, status, err := r.Tasker.CreateRecognizeTask(*file.SaluteId)
	if err != nil {
		return fmt.Errorf("recognizer.CreateRecognizeTask: %w", err)
	}
	file.RecognizeTaskID = &taskId
	file.RecognizeStatus = &status
	return nil
}
