package client

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

const (
	_defaultMaxRetry      = 3
	_defaultRetryInterval = 1000 * time.Millisecond
)

// DeadLetter ...
type DeadLetter interface {
	Save(letter interface{}) error
}

type Client struct {
	httpClient    *http.Client
	host          string
	maxRetry      int
	retryInterval time.Duration
	deadLetter    DeadLetter
	// TODO: rate limiter
}

// New create a client with multiple options or get the default client without providing any options
func New(opts ...Option) *Client {
	cli := &Client{
		httpClient:    &http.Client{},
		maxRetry:      _defaultMaxRetry,
		retryInterval: _defaultRetryInterval,
		deadLetter:    nil,
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
	return c.Parse(ctx, request, response, func(bodyBytes []byte, response interface{}) error {
		return json.Unmarshal(bodyBytes, response)
	})
}

// ParseXML send a request with the given request properties
// Read the body and run xml unmarshaler to fill the given response
func (c *Client) ParseXML(ctx context.Context, request *Request, response interface{}) error {
	return c.Parse(ctx, request, response, func(bodyBytes []byte, response interface{}) error {
		return xml.Unmarshal(bodyBytes, response)
	})
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
func (c *Client) Do(ctx context.Context, request *Request) (*http.Response, error) {
	// TODO: retrying mechanism
	req, err := c.prepareRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
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
