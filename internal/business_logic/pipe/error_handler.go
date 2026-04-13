package pipe

import (
	"context"
	"errors"

	pipeErrors "github.com/a-kuleshov/treplo/internal/business_logic/pipe/errors"
	"github.com/a-kuleshov/treplo/internal/logger"
	"github.com/a-kuleshov/treplo/internal/models"
)

const errorHandlerBufferSize = 100

type errorContainer struct {
	Error error
	File  *models.File
}

type errorTgNotifier interface {
	Notify(replyToMessageId int, chatId int64, message string) error
}

func runErrorHandler(ctx context.Context, tgNotifier errorTgNotifier) (errorChannel chan<- errorContainer) {
	channel := make(chan errorContainer, errorHandlerBufferSize)
	logger := logger.FromContext(ctx)
	go func() {
		for errorContainer := range channel {
			if errorContainer.File == nil {
				logger.Warn("No file in error")
				continue
			}
			userError, ok := errors.AsType[pipeErrors.FileProcessingError](errorContainer.Error)
			message := "В процессе обработки аудио произошла ошибка"
			if ok {
				message = userError.UserMessage
			} else {
				logger.Debug("Not FileProcessingError in channel")
			}
			file := errorContainer.File
			err := tgNotifier.Notify(file.MessageID, file.ChatID, message)
			if err != nil {
				logger.Error("error while user notification", "error", err.Error())
			}
		}
	}()
	return channel
}
