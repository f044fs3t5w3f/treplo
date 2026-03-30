package gigachat

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/a-kuleshov/treplo/pkg/sber/client"
	"github.com/a-kuleshov/treplo/pkg/sber/token"
)

const scopeGigachat = "GIGACHAT_API_PERS"

type GigaChatService struct {
	tokenStorage token.TokenGetter
	wg           *sync.WaitGroup
	client       client.Client
}

func StartGigaChatService(ctx context.Context, clientSecret string) (*GigaChatService, error) {
	tokenStorage, err := token.NewStorage(ctx, clientSecret, scopeGigachat)
	if err != nil {
		return nil, fmt.Errorf("token.NewStorage: %w", err)
	}
	service := GigaChatService{
		tokenStorage: tokenStorage,
		wg:           &sync.WaitGroup{},
		client:       http.DefaultClient,
	}
	return &service, nil
}

func (g *GigaChatService) Stop() {
	g.wg.Wait()
}
