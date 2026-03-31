package client

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	client           Client
	limiter          *rate.Limiter
	retryPolicy      []time.Duration
	rateLimitTimeout time.Duration
}

type ClientOption func(*client)

func WithLimiter(limiter *rate.Limiter, timeout time.Duration) ClientOption {
	return func(c *client) {
		c.limiter = limiter
		c.rateLimitTimeout = timeout
		c.client = &http.Client{}
	}
}

func WithClient(customClient Client) ClientOption {
	return func(c *client) {
		c.client = customClient
	}
}

func WithRetries(retries ...time.Duration) ClientOption {
	return func(c *client) {
		c.retryPolicy = retries
	}
}

func NewClient(options ...ClientOption) *client {
	c := &client{
		client:  &http.Client{},
		limiter: nil,
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func (c *client) wait() error {
	ctx := context.Background()
	if c.rateLimitTimeout != 0 {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), c.rateLimitTimeout)
		ctx = timeoutCtx
		defer cancel()
	}
	return c.limiter.Wait(ctx)
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	if c.limiter != nil {
		if err := c.wait(); err != nil {
			return nil, err
		}
	}
	for _, delay := range c.retryPolicy {
		if resp, err := c.client.Do(req); err != nil {
			time.Sleep(delay)
		} else {
			return resp, nil
		}
	}
	return c.client.Do(req)
}

func (c *client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
