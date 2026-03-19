package waiter

import (
	"context"
	"fmt"
	"time"

	"github.com/a-kuleshov/treplo/internal/models"
)

type statusChecker interface {
	CheckStatus(saluteTaskId string) (string, string, error)
}

type Waiter struct {
	StatusChecker statusChecker
}

// const statusNew = "NEW"
// const statusRunning = "RUNNING"
const statusDone = "DONE"
const statusError = "ERROR"
const statusCanceled = "CANCELED"

func (w *Waiter) Process(ctx context.Context, file *models.File) error {
	if file.RecognizeTaskID == nil {
		return fmt.Errorf("RecognizeTaskID is nil")
	}
	attempts := 3
	for {
		status, responseFileId, err := w.StatusChecker.CheckStatus(*file.RecognizeTaskID)
		if err != nil {
			attempts = attempts - 1
			if attempts == 0 {
				return fmt.Errorf("waiter.CheckStatus: %w", err)
			}
			continue
		} else {
			attempts = 3
			if status == statusDone {
				file.RecognizeStatus = &status
				file.ResponseFileID = &responseFileId
				return nil
			}
			if file.RecognizeStatus == nil || *file.RecognizeStatus != status {
				file.RecognizeStatus = &status
			}
			if status == statusError || status == statusCanceled {
				return fmt.Errorf("waiter.CheckStatus: %s", status)
			}
		}
		time.Sleep(10 * time.Second)
	}
}
