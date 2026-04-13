package pipe

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/content_downloader"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/downloader"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/encoding_detector"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/notifier"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/tasker"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/uploader"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/waiter"
	"github.com/a-kuleshov/treplo/internal/logger"
	"github.com/a-kuleshov/treplo/internal/models"
	"github.com/a-kuleshov/treplo/pkg/sber/salute"
	"github.com/a-kuleshov/treplo/pkg/utils"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/sync/semaphore"
)

const bufSize = 10
const maxOperations = 5

// pipe is the struct containg proccessors of audiofiles.
// It includes the following steps:
// - download file from telegram
// - findout the encoding
// - upload file to Salute speech service
// - create task for speech recognition
// - wait untill task is finished
// - download content of audio
// - send notification
type pipe struct {
	ch chan *models.File
}

func NewPipe(ctx context.Context, repo repository, tgbotapi *tgBotApi.BotAPI, saluteApi *salute.SpeechService, storagePath string) (pipe, error) {
	err := utils.IsDirectoryExistsAndWrible(storagePath)
	if err != nil {
		return pipe{}, fmt.Errorf("utils.IsDirectoryExistsAndWrible: %w, path %s", err, storagePath)
	}
	downloaderProcessor, err := downloader.NewDownloader(tgbotapi.GetFileDirectURL, storagePath)
	if err != nil {
		return pipe{}, fmt.Errorf("downloader.NewDownloader: %w", err)
	}
	tgNotifier := tgNotifier{tgbotapi: tgbotapi}

	processors := []FileProcessor{
		downloaderProcessor,
		&encoding_detector.EncodingDetector{StoragePath: storagePath},
		&uploader.FileUploader{Uploader: saluteApi, StoragePath: storagePath},
		&tasker.Tasker{Tasker: saluteApi},
		&waiter.Waiter{StatusChecker: saluteApi},
		&content_downloader.Tasker{Downloader: saluteApi},
		&notifier.NotifyProccessor{Notifier: tgNotifier},
	}

	channels := make([]chan *models.File, len(processors))
	errorChannel := runErrorHandler(ctx, tgNotifier)
	for i := range channels {
		channels[i] = make(chan *models.File, bufSize)
	}
	pipeInputChannel := channels[0]
	for i, processor := range processors {
		inputChannel := channels[i]
		var outputChannel chan *models.File
		if i < len(channels)-1 {
			outputChannel = channels[i+1]
		}
		processorSemaphore := semaphore.NewWeighted(maxOperations)
		go func() {
			for file := range inputChannel {
				processorSemaphore.Acquire(ctx, 1)
				go func() {
					defer processorSemaphore.Release(1)
					processErr := processor.Process(ctx, file)
					if err != nil {
						logger.FromContext(ctx).Error("File processing error", "fileID", file.ID)
						errorChannel <- errorContainer{processErr, file}
					}
					if err := repo.SaveFile(ctx, file); err != nil {
						logger.FromContext(ctx).Error("pipe: repo.SaveFile", "error", err.Error())
						errorChannel <- errorContainer{err, file}
						return
					}
					if outputChannel != nil && processErr == nil {
						outputChannel <- file
					}
				}()
			}
		}()
	}

	return pipe{
		ch: pipeInputChannel,
	}, nil
}

func (p pipe) Process(ctx context.Context, file *models.File) error {
	p.ch <- file
	return nil // to implement processor interface
}
