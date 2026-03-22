package main

import (
	"context"
	"log/slog"

	"github.com/a-kuleshov/treplo/internal/business_logic/pipe"
	"github.com/a-kuleshov/treplo/internal/db/sql"
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

	repo, err := sql.NewRepository(config.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	tgbotapi, err := tgBotApi.NewBotAPI(config.TgToken)

	if err != nil {
		slog.Error("tgbotapi.NewBotAPI", "error", err)
		panic(err)
	}

	pipe, err := pipe.NewPipe(repo, tgbotapi, speachService, config.StoragePath)
	if err != nil {
		panic(err)
	}

	files, err := repo.ListFiles(context.Background())
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		err := pipe.Process(context.Background(), file)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
