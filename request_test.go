package client_test

import (
	"net/http"
	"testing"

	"github.com/bilginyuksel/client"
	"github.com/stretchr/testify/assert"
)

func TestRequestToURL(t *testing.T) {
	cli := client.New(client.WithHost("http://localhost:3000"))

	testCases := []struct {
		scenario     string
		givenRequest *client.Request
		expectedURL  string
	}{
		{
			scenario:     "GET <base_url/orders/<order-id>?<query-params>",
			givenRequest: cli.NewRequest().Path("/orders/%d", 1231).AddQuery("clientId", "1231321").AddQuery("deviceId", "45555"),
			expectedURL:  "http://localhost:3000/orders/1231?clientId=1231321&deviceId=45555",
		},
		{
			scenario:     "POST <base_url/orders/<order-id>",
			givenRequest: cli.NewRequest().Host("http://0.0.0.0:3000").Method(http.MethodPost).Path("/orders/%d", 1),
			expectedURL:  "http://0.0.0.0:3000/orders/1",
		},
		{
			scenario:     "Add multiple query parameters",
			givenRequest: cli.NewRequest().AddQuery("customerId", "first").AddQuery("customerId", "last"),
			expectedURL:  "http://localhost:3000?customerId=first&customerId=last",
		},
		{
			scenario:     "Add query parameters then override with set",
			givenRequest: cli.NewRequest().AddQuery("customerId", "first").SetQuery("customerId", "last"),
			expectedURL:  "http://localhost:3000?customerId=last",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {

			actualURL, _ := tc.givenRequest.URL()
			assert.Equal(t, tc.expectedURL, actualURL)
		})
	}
}
