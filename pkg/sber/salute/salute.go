package salute

import (
	"context"
	"fmt"
	"sync"

	"github.com/a-kuleshov/treplo/pkg/sber/token"
)

type tokenStorage interface {
	GetToken() (string, error)
}

type SpeachService struct {
	tokenStorage tokenStorage
	wg           *sync.WaitGroup
}

func StartSpeachService(ctx context.Context, clientSecret string) (*SpeachService, error) {
	tokenStorage, err := token.NewStorage(ctx, clientSecret, token.ScopeSaluteSpeech)
	if err != nil {
		return nil, fmt.Errorf("token.NewStorage: %w", err)
	}
	service := SpeachService{
		tokenStorage: tokenStorage,
		wg:           &sync.WaitGroup{},
	}
	return &service, nil
}

func (s *SpeachService) Stop() {
}
