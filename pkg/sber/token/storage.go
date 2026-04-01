package token

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/a-kuleshov/treplo/pkg/client"
	"golang.org/x/time/rate"
)

var ErrNotReady = errors.New("token is not ready")
var ErrWasStoped = errors.New("token getter was stopped")

const delayMargin = 1700 * time.Second

type Storage struct {
	token        string
	lock         *sync.RWMutex
	clientSecret string
	wasStoped    bool
	client       client.Client
}

func NewStorage(ctx context.Context, clientSecret string, scope string) (*Storage, error) {
	token, expiresAt, err := getAccessToken(clientSecret, scope)
	if err != nil {
		return nil, fmt.Errorf("initial getAccessToken: %w", err)
	}
	client := client.NewClient(
		client.WithLimiter(rate.NewLimiter(1, 1), 2*time.Second),
		client.WithRetries(1*time.Second, 2*time.Second, 3*time.Second),
		client.WithClient(&http.Client{Timeout: 100 * time.Second}),
	)
	storage := Storage{
		lock:         &sync.RWMutex{},
		clientSecret: clientSecret,
		token:        token,
		client:       client,
	}

	delta := time.Duration(time.Until(expiresAt))

	go func() {
		for {
			timer := time.NewTimer(delta - delayMargin)
			select {
			case <-ctx.Done():
				slog.Info("Salute speech token service is shutting down")
				timer.Stop()
				storage.lock.Lock()
				storage.wasStoped = true
				storage.lock.Unlock()
				return
			case <-timer.C:
				slog.Info("Get new token")
				token, expiresAt, _ = getAccessToken(clientSecret, scope)
				storage.lock.Lock()
				storage.token = token
				storage.lock.Unlock()
			}
		}
	}()

	return &storage, nil
}

func (s *Storage) GetToken() (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s.wasStoped {
		return "", ErrWasStoped
	}
	return s.token, nil
}
