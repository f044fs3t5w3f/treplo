package salute

import (
	"context"
	"fmt"

	"github.com/a-kuleshov/treplo/pkg/sber/token"
)

type tokenStorage interface {
	GetToken() (string, error)
}

type speachService struct {
	tokenStorage tokenStorage
}

func StartSpeachService(ctx context.Context, clientSecret string) (*speachService, error) {
	tokenStorage, err := token.NewStorage(ctx, clientSecret, token.ScopeSaluteSpeech)
	if err != nil {
		return nil, fmt.Errorf("token.NewStorage: %w", err)
	}
	service := speachService{tokenStorage: tokenStorage}
	return &service, nil
}

func (s *speachService) Stop() {

}
