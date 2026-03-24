package treplo

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/a-kuleshov/treplo/internal/business_logic"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe"
	"github.com/a-kuleshov/treplo/internal/db"
	"github.com/a-kuleshov/treplo/internal/db/sql"
	"github.com/a-kuleshov/treplo/internal/tg"
	"github.com/a-kuleshov/treplo/pkg/sber/gigachat"
	"github.com/a-kuleshov/treplo/pkg/sber/salute"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Stopable interface {
	Stop()
}

type Treplo struct {
	wg       *sync.WaitGroup
	ctx      context.Context
	config   Config
	cancel   func()
	services []Stopable
}

func NewService(config Config) (*Treplo, error) {
	service := &Treplo{
		wg:     &sync.WaitGroup{},
		config: config,
	}
	return service, nil
}

func (t *Treplo) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	t.ctx = ctx
	t.cancel = cancel

	var repository db.Repository
	if t.config.DatabaseDSN != "" {
		var err error
		repository, err = sql.NewRepository(t.config.DatabaseDSN)
		if err != nil {
			slog.Error("Failed to create repository", "error", err)
			panic(err)
		}
	}
	if repository == nil {
		panic("No repository")
	}
	speechService, err := salute.StartSpeachService(ctx, t.config.SaluteSpeechAuthorizationKey)
	if err != nil {
		slog.Error("salute.StartSpeachService", "error", err)
		panic(err)
	}

	gigachatService, err := gigachat.StartGigaChatService(ctx, t.config.GigachatAuthorizationKey)
	if err != nil {
		slog.Error("gigachat.StartGigaChatService", "error", err)
		panic(err)
	}

	tgbotapi, err := tgBotApi.NewBotAPI(t.config.TgToken)

	if err != nil {
		slog.Error("tgbotapi.NewBotAPI", "error", err)
		panic(err)
	}

	// TODO: run add unprocessed files to queue
	fileProcessingPipe, err := pipe.NewPipe(ctx, repository, tgbotapi, speechService, t.config.StoragePath)

	if err != nil {
		slog.Error("pipe.NewPipe", "error", err)
		panic(err)
	}

	business_login := business_logic.NewBusinessLogic(repository, fileProcessingPipe, gigachatService)
	processor := tg.NewProcessor(ctx, business_login, tgbotapi)
	runTGBot(ctx, t.wg, tgbotapi, processor)

	t.services = append(t.services, speechService, gigachatService)
	return nil
}

func (t *Treplo) Stop() error {
	t.cancel()
	for _, service := range t.services {
		service.Stop()
	}
	t.wg.Wait()
	return nil
}

func runTGBot(ctx context.Context, wg *sync.WaitGroup, tgbotapi *tgBotApi.BotAPI, processor *tg.Processor) {
	updateConfig := tgBotApi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := tgbotapi.GetUpdatesChan(updateConfig)
	wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("Stop get tg updates")
				return
			case update := <-updates:
				fmt.Println("in")
				processor.Process(ctx, update)
			}
		}
	})
	slog.Info("runTGBot done")
}
