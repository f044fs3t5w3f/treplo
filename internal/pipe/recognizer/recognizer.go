package recognizer

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/models"
)

type SpeechRecognizer interface {
	CreateRecognizeTask(saluteFileId string) (string, string, error)
}

type Recognizer struct {
	Recognizer SpeechRecognizer
}

func (r *Recognizer) Process(ctx context.Context, file *models.File) error {
	return r.CreateRecognizeTask(ctx, file)
}

func (r *Recognizer) CreateRecognizeTask(ctx context.Context, file *models.File) error {
	taskId, status, err := r.Recognizer.CreateRecognizeTask(*file.SaluteId)
	if err != nil {
		return fmt.Errorf("recognizer.CreateRecognizeTask: %w", err)
	}
	file.RecognizeTaskID = &taskId
	file.RecognizeStatus = &status
	return nil
}
