package client

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const (
	_defaultMaxRetry      = 3
	_defaultRetryInterval = 1000 * time.Millisecond

	_retryIntervalCoef = 1.5
)

type Client struct {
	httpClient    *http.Client
	host          string
	maxRetry      int
	retryInterval time.Duration
	deadLetter    DeadLetter
	rateLimiter   *rate.Limiter
}

// New create a client with multiple options or get the default client without providing any options
func New(opts ...Option) *Client {
	cli := &Client{
		httpClient:    &http.Client{},
		maxRetry:      _defaultMaxRetry,
		retryInterval: _defaultRetryInterval,
	}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
}

// GetJSON execute a get method with the given request and then unmarshal the json response body
func (c *Client) GetJSON(ctx context.Context, request *Request, response interface{}) error {
	return c.ParseJSON(ctx, request.Method(http.MethodGet), response)
}

// GetXML execute a get method with the given request and then unmarshal the xml response body
func (c *Client) GetXML(ctx context.Context, request *Request, response interface{}) error {
	return c.ParseXML(ctx, request.Method(http.MethodGet), response)
}

// ParseJSON send a request with given request properties
// Read the body and run json unmarshaler to fill the given response
func (c *Client) ParseJSON(ctx context.Context, request *Request, response interface{}) error {
	return c.Parse(ctx, request, response, json.Unmarshal)
}

// ParseXML send a request with the given request properties
// Read the body and run xml unmarshaler to fill the given response
func (c *Client) ParseXML(ctx context.Context, request *Request, response interface{}) error {
	return c.Parse(ctx, request, response, xml.Unmarshal)
}

// Parse send a request with the given request properties
// Read the body with the given parser function
func (c *Client) Parse(ctx context.Context, request *Request, response interface{}, parser func(bodyBytes []byte, response interface{}) error) error {
	res, err := c.Do(ctx, request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return parser(responseBytes, response)
}

// Do Execute an http request with the given request
func (c *Client) Do(ctx context.Context, request *Request) (res *http.Response, err error) {
	if err := c.awaitRateLimiter(ctx); err != nil {
		return nil, err
	}

	req, err := c.prepareRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	res, err = c.do(ctx, req, 1)
	if err != nil {
		return nil, err
	}

	// if still 5XX server error then we need to record this request to ensure consistency
	if res.StatusCode >= 500 && res.StatusCode <= 599 {
		if err := c.saveRequest(request, req.URL.String()); err != nil {
			log.Printf("request could not send to deadletter: %v, request: %v\n", err, req)
			return res, fmt.Errorf("letter could not saved: %v", err)
		}
	}

	return res, err
}

func (c *Client) do(ctx context.Context, req *http.Request, retryCount int) (res *http.Response, err error) {
	res, err = c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if c.shouldRetry(retryCount, res.StatusCode) {
		computedRetryInterval := float64(c.retryInterval.Milliseconds()) * math.Pow(_retryIntervalCoef, float64(retryCount))
		time.Sleep(time.Millisecond * time.Duration(computedRetryInterval))
		req.Header.Set("X-Retry", fmt.Sprintf("%d", retryCount))
		return c.do(ctx, req, retryCount+1)
	}

	return res, err
}

func (c *Client) shouldRetry(retryCount int, statusCode int) bool {
	return retryCount <= c.maxRetry && statusCode >= 500 && statusCode <= 599
}

func (c *Client) saveRequest(req *Request, url string) error {
	if c.deadLetter == nil {
		return nil
	}

	return c.deadLetter.Save(&Letter{
		Method:  req.method,
		Body:    req.body,
		Headers: req.headers,
		URL:     url,
	})
}

func (c *Client) prepareRequest(ctx context.Context, request *Request) (*http.Request, error) {
	url, err := request.URL()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, request.method, url, bytes.NewBuffer(request.body))
	if err != nil {
		return nil, err
	}
	req.Header = request.headers

	for _, manipulator := range request.manipulators {
		manipulator(req)
	}

	return req, nil
}

func (c *Client) awaitRateLimiter(ctx context.Context) error {
	if c.rateLimiter == nil {
		return nil
	}
	return c.rateLimiter.Wait(ctx)
}
