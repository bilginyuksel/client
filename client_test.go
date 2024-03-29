package client_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bilginyuksel/client"
	gomock "github.com/golang/mock/gomock"
	"github.com/streetbyters/aduket"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestPostJSON_SuccessfulRequest_ExpectPutMethodInRequest(t *testing.T) {
	s, recorder := aduket.NewServer(http.MethodPost, "/test", aduket.StatusCode(200))
	cli := client.New(client.WithHost(s.URL))
	req := cli.NewRequest().Path("/test").AddHeader("X-R", "req")

	_ = cli.PostJSON(ctx, req, nil)
	// if there is no request captured on PUT aduket.URL/test, it will fail
	recorder.AssertHeaderContains(t, http.Header{"X-R": []string{"req"}})
}

func TestPutJSON_SuccessfulRequest_ExpectPutMethodInRequest(t *testing.T) {
	s, recorder := aduket.NewServer(http.MethodPut, "/test", aduket.StatusCode(200))
	cli := client.New(client.WithHost(s.URL))
	req := cli.NewRequest().Path("/test").AddHeader("X-R", "req")

	_ = cli.PutJSON(ctx, req, nil)
	// if there is no request captured on PUT aduket.URL/test, it will fail
	recorder.AssertHeaderContains(t, http.Header{"X-R": []string{"req"}})
}

func TestGetJSON_SuccessfulRequest_ExpectGETMethodInRequest(t *testing.T) {
	s, recorder := aduket.NewServer(http.MethodGet, "/test", aduket.StatusCode(200))
	cli := client.New(client.WithHost(s.URL), client.WithHTTPClient(&http.Client{}))
	req := cli.NewRequest().Path("/test").
		Method(http.MethodPatch).
		AddHeader("X-R", "req")

	_ = cli.GetJSON(ctx, req, nil)
	// if there is no request captured on GET aduket.URL/test, it will fail
	recorder.AssertHeaderContains(t, http.Header{"X-R": []string{"req"}})
}

func TestGetXML_SuccessfulRequest_ExpectGETMethodInRequest(t *testing.T) {
	s, recorder := aduket.NewServer(http.MethodGet, "/test")
	cli := client.New(client.WithHost(s.URL))
	req := cli.NewRequest().Path("/test").
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
	req := cli.NewRequest().Path("/orders")

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
	req := cli.NewRequest().Path("/test")

	var actualTestStruct Test
	err := cli.ParseXML(ctx, req, &actualTestStruct)

	assert.Nil(t, err)
	assert.Equal(t, expectedTestStruct, &actualTestStruct)
}

func TestParse_CorruptedBody_ReturnErr(t *testing.T) {
	s, _ := aduket.NewServer(http.MethodGet, "/orders", aduket.CorruptedBody())
	cli := client.New(client.WithHost(s.URL))
	req := cli.NewRequest().Path("/orders")

	err := cli.Parse(ctx, req, nil, nil)
	assert.NotNil(t, err)
}

func TestParse_DoFailed_ReturnErr(t *testing.T) {
	cli := client.New(client.WithHost("local:host:3000"))
	req := cli.NewRequest() // to fail new request with context

	err := cli.Parse(ctx, req, nil, nil)
	assert.NotNil(t, err)
}

func TestDo_5XXStatusCode_RetryNTimes(t *testing.T) {
	var (
		_maxRetryCount = 3
		_retryInterval = 1 * time.Millisecond

		count             int
		xRetryHeaderValue string
		expectedCount     = 4
	)

	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		count++
		xRetryHeaderValue = r.Header.Get("X-Retry")
		rw.WriteHeader(500)
	}))

	cli := client.New(client.WithHost(s.URL), client.WithRetry(_maxRetryCount, _retryInterval))
	_, _ = cli.Do(ctx, cli.NewRequest())

	assert.Equal(t, expectedCount, count)
	assert.Equal(t, "3", xRetryHeaderValue)
}

func TestDo_5XXStatusCode_RetryIntervalShouldIncreaseExponentialy(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}
	var (
		_maxRetryCount = 2
		_retryInterval = 500 * time.Millisecond

		capturedTimes []time.Time
	)

	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		capturedTimes = append(capturedTimes, time.Now())
		rw.WriteHeader(500)
	}))

	cli := client.New(client.WithHost(s.URL), client.WithRetry(_maxRetryCount, _retryInterval))
	_, _ = cli.Do(ctx, cli.NewRequest())

	assert.GreaterOrEqual(t, capturedTimes[1].Sub(capturedTimes[0]), 500*time.Millisecond)
	assert.GreaterOrEqual(t, capturedTimes[2].Sub(capturedTimes[1]), 750*time.Millisecond)
}

func TestDo_ReachMaxRetry_SaveRequestToDeadLetter(t *testing.T) {
	s, _ := aduket.NewServer(http.MethodGet, "/test", aduket.StatusCode(502))

	mockDeadLetter := client.NewMockDeadLetter(gomock.NewController(t))
	mockDeadLetter.EXPECT().Save(&client.Letter{
		Method:  "GET",
		URL:     fmt.Sprintf("%s/test", s.URL),
		Headers: map[string][]string{"X-Retry": {"3"}},
	})

	cli := client.New(client.WithHost(s.URL), client.WithDeadLetter(mockDeadLetter), client.WithRetry(3, 1*time.Millisecond))
	_, _ = cli.Do(ctx, cli.NewRequest().Path("/test"))
}

func TestDo_SaveDeadLetterFailedAfterReachingMaxRetry_ReturnErr(t *testing.T) {
	s, _ := aduket.NewServer(http.MethodGet, "/test", aduket.StatusCode(502))

	mockDeadLetter := client.NewMockDeadLetter(gomock.NewController(t))
	mockDeadLetter.EXPECT().Save(gomock.Any()).Return(errors.New("error"))

	cli := client.New(client.WithHost(s.URL), client.WithDeadLetter(mockDeadLetter), client.WithRetry(1, 1*time.Millisecond))
	_, err := cli.Do(ctx, cli.NewRequest().Path("/test"))

	assert.NotNil(t, err)
	assert.Equal(t, "letter could not saved: error", err.Error())
}

func TestDo_MultipleRequestsAtOnce_RateLimitAndWait(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	startTime := time.Now()
	server, _ := aduket.NewServer(http.MethodGet, "/test", aduket.StatusCode(200))
	cli := client.New(client.WithHost(server.URL), client.WithRateLimit(65*time.Millisecond, 50))

	for i := 0; i < 100; i++ {
		ctx = context.Background()
		_, _ = cli.Do(ctx, cli.NewRequest().Path("/test"))
	}

	assert.GreaterOrEqual(t, time.Since(startTime), 3*time.Second)
}

func TestDo_ValidRequest_CaptureTheHTTPRequestWithTheGivenParameters(t *testing.T) {
	s, recorder := aduket.NewServer(http.MethodPost, "/basket/1")
	cli := client.New(client.WithHost(s.URL))

	req := cli.NewRequest().
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
	recorder.AssertHeaderContains(t, expectedHeaders)
	recorder.AssertStringBodyEqual(t, "hello world")
	recorder.AssertQueryParamEqual(t, "customerId", []string{"12321"})
}

func TestDo_WrongHTTPMethod_ReturnErr(t *testing.T) {
	cli := client.New(client.WithHost("http://localhost:3000"))
	req := cli.NewRequest().Method(",") // to fail new request with context

	_, err := cli.Do(ctx, req)

	assert.NotNil(t, err)
}

func TestDo_InvalidHost_ReturnErr(t *testing.T) {
	cli := client.New(client.WithHost("local:host:3000"))
	req := cli.NewRequest() // to fail new request with context

	_, err := cli.Do(ctx, req)

	assert.NotNil(t, err)
}
