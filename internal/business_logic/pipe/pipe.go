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
	"github.com/a-kuleshov/treplo/internal/models"
	"github.com/a-kuleshov/treplo/pkg/sber/salute"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type pipe struct {
	processors []FileProcessor
	repo       repository
}

type tgNotifier struct {
	tgbotapi *tgBotApi.BotAPI
}

func (tn tgNotifier) Notify(replyToMessageId int, chatId int64, message string) error {
	tgMessage := tgBotApi.NewMessage(chatId, message)
	tgMessage.ReplyToMessageID = replyToMessageId
	_, err := tn.tgbotapi.Send(tgMessage)
	return err
}

func NewPipe(repo repository, tgbotapi *tgBotApi.BotAPI, saluteApi *salute.SpeachService, storagePath string) (pipe, error) {
	downloaderProcessor, err := downloader.NewDownloader(tgbotapi.GetFileDirectURL, storagePath)
	if err != nil {
		return pipe{}, fmt.Errorf("downloader.NewDownloader: %w", err)
	}
	return pipe{
		processors: []FileProcessor{
			downloaderProcessor,
			&uploader.FileUploader{Uploader: saluteApi},
			&tasker.Tasker{Tasker: saluteApi},
			&waiter.Waiter{StatusChecker: saluteApi},
			&content_downloader.Tasker{Downloader: saluteApi},
			&notifier.NotifyProccessor{Notifier: tgNotifier{tgbotapi: tgbotapi}},
		},
		repo: repo,
	}, nil
}

func (p pipe) Process(ctx context.Context, file *models.File) error {
	for _, processor := range p.processors {
		if err := processor.Process(ctx, file); err != nil {
			return err
		}
		if err := p.repo.SaveFile(ctx, file); err != nil {
			return fmt.Errorf("repo.SaveFile: %w", err)
		}
	}
	return nil
}
