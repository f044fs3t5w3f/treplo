package treplo

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/a-kuleshov/treplo/internal/business_logic"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe"
	"github.com/a-kuleshov/treplo/internal/db"
	"github.com/a-kuleshov/treplo/internal/db/sql"
	"github.com/a-kuleshov/treplo/internal/logger"
	"github.com/a-kuleshov/treplo/internal/tg"
	"github.com/a-kuleshov/treplo/pkg/sber/gigachat"
	"github.com/a-kuleshov/treplo/pkg/sber/salute"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Stopable interface {
	Stop()
}

// Treplo struct is the main service of the application that manage all the subservices: database, telegram bot handeler, salute speech and gigachat
type Treplo struct {
	wg       *sync.WaitGroup
	config   Config
	cancel   func()
	services []Stopable
}

func NewService(config Config) *Treplo {
	service := &Treplo{
		wg:     &sync.WaitGroup{},
		config: config,
	}
	return service
}

func (t *Treplo) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	var repository db.Repository
	if t.config.DatabaseDSN != "" {
		var err error
		repository, err = sql.NewRepository(t.config.DatabaseDSN)
		if err != nil {
			logger.FromContext(ctx).Error("Failed to create repository", "error", err)
			return err
		}
	}
	if repository == nil {
		logger.FromContext(ctx).Error("No repository")
		return errors.New("no repository")
	}
	speechService, err := salute.StartSpeechService(ctx, t.config.SaluteSpeechAuthorizationKey)
	if err != nil {
		logger.FromContext(ctx).Error("salute.StartSpeechService", "error", err)
		return err
	}

	gigachatService, err := gigachat.StartGigaChatService(ctx, t.config.GigachatAuthorizationKey)
	if err != nil {
		logger.FromContext(ctx).Error("gigachat.StartGigaChatService", "error", err)
		return err
	}

	tgbotapi, err := tgBotApi.NewBotAPI(t.config.TgToken)
	if err != nil {
		logger.FromContext(ctx).Error("tgbotapi.NewBotAPI", "error", err)
		return err
	}

	fileProcessingPipe, err := pipe.NewPipe(ctx, repository, tgbotapi, speechService, t.config.StoragePath)
	if err != nil {
		logger.FromContext(ctx).Error("pipe.NewPipe", "error", err)
		return err
	}

	go runUnprocessedFilesProcessing(ctx, repository, fileProcessingPipe)

	business_login := business_logic.NewBusinessLogic(repository, fileProcessingPipe, gigachatService)
	processor := tg.NewProcessor(ctx, business_login, tgbotapi)
	runTGBot(ctx, t.wg, tgbotapi, processor)

	t.services = append(t.services, speechService, gigachatService)
	return nil
}

func (t *Treplo) Stop() {
	t.cancel()
	timer := time.AfterFunc(3*time.Second, func() { os.Exit(130) })
	defer timer.Stop()
	for _, service := range t.services {
		service.Stop()
	}
	t.wg.Wait()
}

func runTGBot(ctx context.Context, wg *sync.WaitGroup, tgbotapi *tgBotApi.BotAPI, processor *tg.Processor) {
	updateConfig := tgBotApi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := tgbotapi.GetUpdatesChan(updateConfig)
	wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				logger.FromContext(ctx).Info("Stop get tg updates")
				return
			case update := <-updates:
				processor.Process(ctx, update)
			}
		}
	})
	logger.FromContext(ctx).Info("runTGBot done")
}

func runUnprocessedFilesProcessing(ctx context.Context, repository db.Repository, fileProcessingPipe pipe.FileProcessor) {
	unprocessedFiles, err := repository.ListNewFiles(ctx)
	if err == nil {
		for _, file := range unprocessedFiles {
			select {
			case <-ctx.Done():
				return
			default:
			}
			fileProcessingPipe.Process(ctx, file)
			time.Sleep(1 * time.Second) // to avoid pipline overloading
		}
	} else {
		logger.FromContext(ctx).Warn("failed to select unprocessed files", "error", err.Error())
	}
}
