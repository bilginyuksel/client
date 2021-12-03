package client

import (
	"context"
	"fmt"
	"net/http"

	urlpkg "net/url"
)

// Request use this request struct to send http requests
type Request struct {
	ctx     context.Context
	method  string
	host    string
	path    string
	body    []byte
	query   map[string][]string
	headers http.Header

	manipulators []func(r *http.Request)
}

// NewRequest creates a new request with the given context
func (c *Client) NewRequest(ctx context.Context) *Request {
	return &Request{
		ctx:     ctx,
		host:    c.host,
		method:  http.MethodGet,
		query:   make(map[string][]string),
		headers: make(map[string][]string),
	}
}

// Host set the host
func (r *Request) Host(host string) *Request {
	r.host = host
	return r
}

// Method set a method to given request
func (r *Request) Method(method string) *Request {
	r.method = method
	return r
}

// Path sets the given request path to request.
// To avoid fmt.Sprintf call while calling the function you can
// directly give the formatted string as path variable and options after that
// function will automatically insert the paremeters
func (r *Request) Path(path string, opts ...interface{}) *Request {
	r.path = fmt.Sprintf(path, opts...)
	return r
}

// SetQuery if given query key currently has a value it will replace it with the given value
// if the key does not exists it will add a new key value
func (r *Request) SetQuery(key string, value ...string) *Request {
	r.query[key] = value
	return r
}

// AddQuery if given query key currently has a value it will add another one
// if the key does not exists it will add a new key value
func (r *Request) AddQuery(key string, value ...string) *Request {
	r.query[key] = append(r.query[key], value...)
	return r
}

// SetHeader if given header key currently has a value it will replace it with the given value
// if the key does not exists it will add a new key value
func (r *Request) SetHeader(key string, value ...string) *Request {
	r.headers[key] = value
	return r
}

// AddHeader if given header key currently has a value it will add another one
// if the key does not exists it will add a new key value
func (r *Request) AddHeader(key string, value ...string) *Request {
	r.headers[key] = append(r.headers[key], value...)
	return r
}

func (r *Request) SetBasicAuth(username, password string) *Request {
	r.manipulators = append(r.manipulators, func(r *http.Request) {
		r.SetBasicAuth(username, password)
	})
	return r
}

func (r *Request) URL() (string, error) {
	rawpath := fmt.Sprintf("%s%s", r.host, r.path)
	url, err := urlpkg.Parse(rawpath)
	if err != nil {
		return "", err
	}

	url.RawQuery = r.getEncodedQueryParameters()
	return url.String(), nil
}

func (r *Request) getEncodedQueryParameters() string {
	queryBuilder := urlpkg.Values{}
	for key, values := range r.query {
		for _, v := range values {
			queryBuilder.Add(key, v)
		}
	}
	return queryBuilder.Encode()
}
