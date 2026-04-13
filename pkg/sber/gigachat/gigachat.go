package gigachat

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/a-kuleshov/treplo/pkg/client"
	"github.com/a-kuleshov/treplo/pkg/sber/token"
	"golang.org/x/time/rate"
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
	client := client.NewClient(
		client.WithLimiter(rate.NewLimiter(100, 100), 2*time.Second),
		client.WithRetries(1*time.Second, 2*time.Second, 3*time.Second),
		client.WithClient(&http.Client{Timeout: 10 * time.Second}),
	)
	service := GigaChatService{
		tokenStorage: tokenStorage,
		wg:           &sync.WaitGroup{},
		client:       client,
	}
	return &service, nil
}

func (g *GigaChatService) Stop() {
	g.wg.Wait()
}
