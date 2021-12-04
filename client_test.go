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

var ctx = context.Background()

func TestGetJSON_SuccessfulRequest_ExpectGETMethodInRequest(t *testing.T) {
	s, recorder := aduket.NewServer(http.MethodGet, "/test")
	cli := client.New(client.WithHost(s.URL))
	req := cli.NewRequest(ctx).Path("/test").
		Method(http.MethodPatch).
		AddHeader("X-R", "req")

	_ = cli.GetJSON(ctx, req, nil)
	// if there is no request captured on GET aduket.URL/test, it will fail
	recorder.AssertHeaderContains(t, http.Header{"X-R": []string{"req"}})
}

func TestGetXML_SuccessfulRequest_ExpectGETMethodInRequest(t *testing.T) {
	s, recorder := aduket.NewServer(http.MethodGet, "/test")
	cli := client.New(client.WithHost(s.URL))
	req := cli.NewRequest(ctx).Path("/test").
		Method(http.MethodPatch).
		AddHeader("X-R", "req")

	_ = cli.GetXML(ctx, req, nil)
	// if there is no request captured on GET aduket.URL/test, it will fail
	recorder.AssertHeaderContains(t, http.Header{"X-R": []string{"req"}})
}

func TestParseJSON_SuccessfulRequest_FillGivenResponseStruct(t *testing.T) {
	type Test struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}
	expectedTestStruct := &Test{
		Firstname: "firstname",
		Lastname:  "lastname",
	}
	s, _ := aduket.NewServer(http.MethodGet, "/orders", aduket.JSONBody(expectedTestStruct))
	cli := client.New(client.WithHost(s.URL))
	req := cli.NewRequest(ctx).Path("/orders")

	var actualTestStruct Test
	err := cli.ParseJSON(ctx, req, &actualTestStruct)

	assert.Nil(t, err)
	assert.Equal(t, expectedTestStruct, &actualTestStruct)
}

func TestParseXML_SuccessfulRequest_FillGivenResponseStruct(t *testing.T) {
	type Test struct {
		Firstname string `xml:"firstname"`
		Lastname  string `xml:"lastname"`
	}
	expectedTestStruct := &Test{
		Firstname: "firstname",
		Lastname:  "lastname",
	}
	s, _ := aduket.NewServer(http.MethodGet, "/test", aduket.XMLBody(expectedTestStruct))
	cli := client.New(client.WithHost(s.URL))
	req := cli.NewRequest(ctx).Path("/test")

	var actualTestStruct Test
	err := cli.ParseXML(ctx, req, &actualTestStruct)

	assert.Nil(t, err)
	assert.Equal(t, expectedTestStruct, &actualTestStruct)
}

func TestParse_CorruptedBody_ReturnErr(t *testing.T) {
	s, _ := aduket.NewServer(http.MethodGet, "/orders", aduket.CorruptedBody())
	cli := client.New(client.WithHost(s.URL))
	req := cli.NewRequest(ctx).Path("/orders")

	err := cli.Parse(ctx, req, nil, nil)
	assert.NotNil(t, err)
}

func TestParse_DoFailed_ReturnErr(t *testing.T) {
	cli := client.New(client.WithHost("local:host:3000"))
	req := cli.NewRequest(ctx) // to fail new request with context

	err := cli.Parse(ctx, req, nil, nil)
	assert.NotNil(t, err)
}

func TestDo_ValidRequest_CaptureTheHTTPRequestWithTheGivenParameters(t *testing.T) {
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
	cli := client.New(client.WithHost("http://localhost:3000"))
	req := cli.NewRequest(ctx).Method(",") // to fail new request with context

	_, err := cli.Do(ctx, req)

	assert.NotNil(t, err)
}

func TestDo_InvalidHost_ReturnErr(t *testing.T) {
	cli := client.New(client.WithHost("local:host:3000"))
	req := cli.NewRequest(ctx) // to fail new request with context

	_, err := cli.Do(ctx, req)

	assert.NotNil(t, err)
}
