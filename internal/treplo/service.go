package treplo

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/a-kuleshov/treplo/internal/business_logic"
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe"
	"github.com/a-kuleshov/treplo/internal/db"
	"github.com/a-kuleshov/treplo/internal/tg"
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
		repository, err = db.NewRepository(t.config.DatabaseDSN)
		if err != nil {
			slog.Error("Failed to create repository", "error", err)
			panic(err)
		}
	}
	if repository == nil {
		panic("No repository")
	}

	tgbotapi, err := tgBotApi.NewBotAPI(t.config.TgToken)

	if err != nil {
		slog.Error("tgbotapi.NewBotAPI", "error", err)
		panic(err)
	}
	speechService, err := salute.StartSpeachService(ctx, t.config.SaluteSpeechAuthorizationKey)
	if err != nil {
		slog.Error("salute.StartSpeachService", "error", err)
		panic(err)
	}

	processors, err := pipe.NewPipe(repository, tgbotapi, speechService, t.config.StoragePath)
	if err != nil {
		slog.Error("pipe.NewPipe", "error", err)
		panic(err)
	}

	mchncs := business_logic.NewBusinessLogic(repository, processors)
	processor := tg.NewProcessor(ctx, mchncs, tgbotapi)
	runTGBot(ctx, t.wg, tgbotapi, processor)

	// saluteSpeachService, err := salute.StartSpeachService(ctx, t.config.SaluteSpeechAuthorizationKey)
	// if err != nil {
	// 	return err
	// }
	// t.services = append(t.services, saluteSpeachService)
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
