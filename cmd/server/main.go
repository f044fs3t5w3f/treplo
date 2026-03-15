package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/a-kuleshov/treplo/internal/service"
	"github.com/a-kuleshov/treplo/pkg/configuration"
)

func main() {
	var config service.Config
	configuration.ScanConfig(&config, nil)
	fmt.Println(config)

	treploSerivce, err := service.NewService(config)
	if err != nil {
		slog.Error(err.Error())
	}
	treploSerivce.Run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGQUIT)
	<-signals
	fmt.Println("stopping")
	treploSerivce.Stop()
}
