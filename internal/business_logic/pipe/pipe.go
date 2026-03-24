package pipe

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/content_downloader"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/downloader"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/notifier"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/tasker"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/uploader"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe/waiter"
	"github.com/a-kuleshov/treplo/internal/logger"
	"github.com/a-kuleshov/treplo/internal/models"
	"github.com/a-kuleshov/treplo/pkg/sber/salute"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const bufSize = 10

type pipe struct {
	// processors []FileProcessor
	// repo       repository
	ch chan *models.File
}

func NewPipe(ctx context.Context, repo repository, tgbotapi *tgBotApi.BotAPI, saluteApi *salute.SpeachService, storagePath string) (pipe, error) {
	downloaderProcessor, err := downloader.NewDownloader(tgbotapi.GetFileDirectURL, storagePath)
	if err != nil {
		return pipe{}, fmt.Errorf("downloader.NewDownloader: %w", err)
	}
	// TODO: make pools of processors instead of single one
	tgNotifier := tgNotifier{tgbotapi: tgbotapi}
	processors := []FileProcessor{
		downloaderProcessor,
		&uploader.FileUploader{Uploader: saluteApi},
		&tasker.Tasker{Tasker: saluteApi},
		&waiter.Waiter{StatusChecker: saluteApi},
		&content_downloader.Tasker{Downloader: saluteApi},
		&notifier.NotifyProccessor{Notifier: tgNotifier},
	}
	channels := make([]chan *models.File, len(processors))

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
		go func() {
			for file := range inputChannel {
				processor.Process(ctx, file)
				if err := repo.SaveFile(ctx, file); err != nil {
					logger.FromContext(ctx).Error("pipe: repo.SaveFile", "error", err.Error())
					tgNotifier.Notify(file.MessageID, file.ChatID, "В процессе обработки аудио произошла ошибка")
					continue
				}
				if outputChannel != nil {
					outputChannel <- file
				}
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
