package token

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

var ErrNotReady = errors.New("token is not ready")

type Storage struct {
	token        string
	lock         *sync.RWMutex
	clientSecret string
}

func NewStorage(ctx context.Context, clientSecret string) (*Storage, error) {
	token, expiresAt, err := getAccessToken(clientSecret)
	if err != nil {
		return nil, fmt.Errorf("initial getAccessToken: %w", err)
	}

	storage := Storage{
		lock:         &sync.RWMutex{},
		clientSecret: clientSecret,
		token:        token,
	}

	delta := time.Duration(time.Until(expiresAt).Seconds())
	fmt.Println(delta)

	go func() {
		for {
			timer := time.NewTimer((delta - 1794) * time.Second)
			select {
			case <-ctx.Done():
				slog.Info("Salute speach token service is shutting down")
				return
			case <-timer.C:
				slog.Info("Get new token")
				token, expiresAt, _ = getAccessToken(clientSecret)
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
	return s.token, nil
}
