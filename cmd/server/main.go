package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/a-kuleshov/treplo/internal/treplo"
	"github.com/a-kuleshov/treplo/pkg/configuration"
)

func main() {
	var config treplo.Config
	configuration.ScanConfig(&config, nil)

	treploSerivce := treplo.NewService(config)
	err := treploSerivce.Run()
	if err != nil {
		slog.Error(err.Error())
		panic(err.Error())
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGQUIT)
	<-signals
	fmt.Println("stopping")
	treploSerivce.Stop()
}
