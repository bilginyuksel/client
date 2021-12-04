package main

import (
	"context"

	"github.com/bilginyuksel/client"
	"github.com/yudai/pp"
)

type MockClient struct {
	*client.Client
}

func NewMockClient(options client.Option) *MockClient {
	return &MockClient{
		Client: client.New(options),
	}
}

type Coffee struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
}

func main() {
	c := NewMockClient(client.WithHost("https://api.sampleapis.com"))
	req := c.NewRequest(context.Background()).Path("/coffee/hot")

	var hotCoffees []Coffee
	if err := c.GetJSON(context.Background(), req, &hotCoffees); err != nil {
		panic(err)
	}
	pp.Println(hotCoffees)

	var icedCoffees []Coffee
	if err := c.GetJSON(context.Background(), req, &icedCoffees); err != nil {
		panic(err)
	}
	pp.Println(icedCoffees)
}
