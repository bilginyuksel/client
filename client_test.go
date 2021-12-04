package client_test

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/bilginyuksel/client"
	"github.com/streetbyters/aduket"
	"github.com/stretchr/testify/assert"
)

func TestDo_ValidRequest_CaptureTheHTTPRequestWithTheGivenParameters(t *testing.T) {
	ctx := context.Background()
	s, recorder := aduket.NewServer(http.MethodPost, "/basket/1")
	cli := client.New(client.WithHost(s.URL))

	req := cli.NewRequest(ctx).
		Method(http.MethodPost).
		Path("/basket/%d", 1).
		AddQuery("customerId", "12321").
		AddHeader("X-Redirect", "12.13.22.100").
		SetHeader("X-Redirect", "0.0.0.0", "8.8.8.8").
		AddHeader("User-Agent", "bilginyuksel").
		SetBasicAuth("test", "test").
		Body([]byte("hello world"))

	_, err := cli.Do(ctx, req)

	expectedHeaders := http.Header{
		"X-Redirect":    []string{"0.0.0.0", "8.8.8.8"},
		"User-Agent":    []string{"bilginyuksel"},
		"Authorization": []string{"Basic dGVzdDp0ZXN0"},
	}

	assert.Nil(t, err)
	log.Println(recorder.Header)
	recorder.AssertHeaderContains(t, expectedHeaders)
	recorder.AssertStringBodyEqual(t, "hello world")
	recorder.AssertQueryParamEqual(t, "customerId", []string{"12321"})
}

func TestDo_WrongHTTPMethod_ReturnErr(t *testing.T) {
	ctx := context.Background()
	cli := client.New(client.WithHost("http://localhost:3000"))
	req := cli.NewRequest(ctx).Method(",") // to fail new request with context

	_, err := cli.Do(ctx, req)

	assert.NotNil(t, err)
}

func TestDo_InvalidHost_ReturnErr(t *testing.T) {
	ctx := context.Background()
	cli := client.New(client.WithHost("local:host:3000"))
	req := cli.NewRequest(ctx) // to fail new request with context

	_, err := cli.Do(ctx, req)

	assert.NotNil(t, err)
}
