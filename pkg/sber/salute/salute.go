package salute

import (
	"context"
	"fmt"
	"sync"

	"github.com/a-kuleshov/treplo/pkg/sber/token"
)

const ScopeSaluteSpeech = "SALUTE_SPEECH_PERS"

type SpeechService struct {
	tokenStorage token.TokenGetter
	wg           *sync.WaitGroup
}

// TODO: reuse client intead of creating a new one for every request

func StartSpeechService(ctx context.Context, clientSecret string) (*SpeechService, error) {
	tokenStorage, err := token.NewStorage(ctx, clientSecret, ScopeSaluteSpeech)
	if err != nil {
		return nil, fmt.Errorf("token.NewStorage: %w", err)
	}
	service := SpeechService{
		tokenStorage: tokenStorage,
		wg:           &sync.WaitGroup{},
	}
	return &service, nil
}

func (s *SpeechService) Stop() {
	s.wg.Wait()
}
