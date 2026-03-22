package gigachat

import (
	"context"
	"fmt"
	"sync"

	"github.com/a-kuleshov/treplo/pkg/sber/token"
)

const scopeGigachat = "GIGACHAT_API_PERS"

type GigaChatService struct {
	tokenStorage token.TokenGetter
	wg           *sync.WaitGroup
}

func StartGigaChatService(ctx context.Context, clientSecret string) (*GigaChatService, error) {
	tokenStorage, err := token.NewStorage(ctx, clientSecret, scopeGigachat)
	if err != nil {
		return nil, fmt.Errorf("token.NewStorage: %w", err)
	}
	service := GigaChatService{
		tokenStorage: tokenStorage,
		wg:           &sync.WaitGroup{},
	}
	return &service, nil
}

func (g *GigaChatService) Stop() {
	g.wg.Wait()
}
