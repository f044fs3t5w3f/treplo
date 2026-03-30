package salute

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/a-kuleshov/treplo/pkg/sber/client"
	"github.com/a-kuleshov/treplo/pkg/sber/token"
)

const ScopeSaluteSpeech = "SALUTE_SPEECH_PERS"

type SpeechService struct {
	tokenStorage token.TokenGetter
	wg           *sync.WaitGroup
	client       client.Client
}

func StartSpeechService(ctx context.Context, clientSecret string) (*SpeechService, error) {
	tokenStorage, err := token.NewStorage(ctx, clientSecret, ScopeSaluteSpeech)
	if err != nil {
		return nil, fmt.Errorf("token.NewStorage: %w", err)
	}
	// TODO: add ratelimit
	client := &http.Client{Timeout: 10 * time.Second}
	service := SpeechService{
		tokenStorage: tokenStorage,
		wg:           &sync.WaitGroup{},
		client:       client,
	}
	return &service, nil
}

func (s *SpeechService) Stop() {
	s.wg.Wait()
}
