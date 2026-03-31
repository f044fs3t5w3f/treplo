package client

import (
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

type failedClient struct{}

func (c failedClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("failed")
}

func (c failedClient) Get(string) (*http.Response, error) {
	return nil, errors.New("failed")
}

type LongClient struct {
	time time.Duration
}

func (lc LongClient) Do(req *http.Request) (*http.Response, error) {
	time.Sleep(lc.time)
	return &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
	}, nil
}

func (lc LongClient) Get(string) (*http.Response, error) {
	time.Sleep(lc.time)
	return &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
	}, nil
}
func TestClient_retries(t *testing.T) {
	testCases := []struct {
		name      string
		durations []time.Duration
	}{
		{"One", []time.Duration{time.Second}},
		{"two", []time.Duration{500 * time.Millisecond, time.Second}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := NewClient(
				WithClient(failedClient{}),
				WithRetries(testCase.durations...),
			)
			req := &http.Request{}
			start := time.Now()
			c.Do(req)
			duration := time.Since(start)
			assert.GreaterOrEqual(t, duration, sum(testCase.durations))
			assert.GreaterOrEqual(t, duration, sum(testCase.durations))
		})
	}
}

func TestClient_rateLimit_timeout(t *testing.T) {
	c := NewClient(
		WithLimiter(rate.NewLimiter(1, 1), 500*time.Millisecond),
		WithClient(LongClient{0}),
	)
	req := &http.Request{}
	wg := sync.WaitGroup{}
	wg.Go(
		func() {
			c.Do(req)
		},
	) // reduce limiter
	wg.Wait()
	_, err := c.Do(req)
	assert.Error(t, err)
}
func TestClient_rateLimit_wait(t *testing.T) {
	c := NewClient(
		WithLimiter(rate.NewLimiter(1, 1), 2500*time.Millisecond),
		WithClient(LongClient{0}),
	)
	req := &http.Request{}
	wg := sync.WaitGroup{}
	wg.Go(
		func() {
			c.Do(req)
		},
	) // reduce limiter
	start := time.Now()
	wg.Wait()
	_, err := c.Do(req)
	duration := time.Since(start)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, duration, 1*time.Second)
}

func sum(durations []time.Duration) time.Duration {
	total := 0 * time.Second
	for _, duration := range durations {
		total += duration
	}
	return total
}
