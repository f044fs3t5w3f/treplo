package service

import (
	"context"
	"sync"

	"github.com/a-kuleshov/treplo/pkg/sber/salute"
)

type Stopable interface {
	Stop()
}

type Service struct {
	wg       *sync.WaitGroup
	ctx      context.Context
	config   Config
	cancel   func()
	services []Stopable
}

func NewService(config Config) (*Service, error) {
	service := &Service{
		wg:     &sync.WaitGroup{},
		config: config,
	}
	return service, nil
}

func (s *Service) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	saluteSpeachService, err := salute.StartSpeachService(ctx, s.config.SaluteSpeechAuthorizationKey)
	if err != nil {
		return err
	}
	s.services = append(s.services, saluteSpeachService)
	return nil
}

func (s *Service) Stop() error {
	s.cancel()
	for _, service := range s.services {
		service.Stop()
	}
	s.wg.Wait()
	return nil
}
