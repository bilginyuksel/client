package client

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// Option is a function that configures a client
type Option func(c *Client)

// WithHost create client option function with host
func WithHost(host string) Option {
	return func(c *Client) {
		c.host = host
	}
}

// WithRetry create client option function with retrying properties
func WithRetry(maxRetry int, retryInterval time.Duration) Option {
	return func(c *Client) {
		c.maxRetry = maxRetry
		c.retryInterval = retryInterval
	}
}

// WithDeadLetter create client option function with deadletter properties
func WithDeadLetter(deadLetter DeadLetter) Option {
	return func(c *Client) {
		c.deadLetter = deadLetter
	}
}

// WithHTTPClient create client option function with http client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithRateLimit create client option function with http client
func WithRateLimit(interval time.Duration, requests int) Option {
	return func(c *Client) {
		c.rateLimiter = rate.NewLimiter(rate.Every(interval), requests)
	}
}
