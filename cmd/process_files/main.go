package main

import (
	"context"
	"log/slog"

	"github.com/a-kuleshov/treplo/internal/db"
	"github.com/a-kuleshov/treplo/internal/pipe"
	"github.com/a-kuleshov/treplo/internal/treplo"
	"github.com/a-kuleshov/treplo/pkg/configuration"
	"github.com/a-kuleshov/treplo/pkg/sber/salute"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	var config treplo.Config
	configuration.ScanConfig(&config, nil)
	speachService, err := salute.StartSpeachService(context.Background(), config.SaluteSpeechAuthorizationKey)
	if err != nil {
		panic(err)
	}

	repo, err := db.NewRepository(config.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	tgbotapi, err := tgBotApi.NewBotAPI(config.TgToken)

	if err != nil {
		slog.Error("tgbotapi.NewBotAPI", "error", err)
		panic(err)
	}

	pipe := pipe.NewPipe(repo, tgbotapi.GetFileDirectURL, speachService)

	files, err := repo.ListFiles(context.Background())
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.ID != 1 {
			continue
		}
		pipe.Process(context.Background(), &file)
	}
}
