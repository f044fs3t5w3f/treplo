package salute

import (
	"context"
	"fmt"
	"sync"

	"github.com/a-kuleshov/treplo/pkg/sber/token"
)

const ScopeSaluteSpeech = "SALUTE_SPEECH_PERS"

type SpeachService struct {
	tokenStorage token.TokenGetter
	wg           *sync.WaitGroup
}

// TODO: reuse client intead of creating a new one for every request

func StartSpeachService(ctx context.Context, clientSecret string) (*SpeachService, error) {
	tokenStorage, err := token.NewStorage(ctx, clientSecret, ScopeSaluteSpeech)
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
	s.wg.Wait()
}
