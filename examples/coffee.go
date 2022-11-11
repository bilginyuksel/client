package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bilginyuksel/client"
	"github.com/labstack/echo/v4"
)

type MockClient struct {
	*client.Client
}

func NewMockClient(options ...client.Option) *MockClient {
	return &MockClient{
		Client: client.New(options...),
	}
}

type Coffee struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
}

func main() {
	c := NewMockClient(client.WithHost("http://localhost:8080"),
		client.WithRetry(3, time.Millisecond*200),
	)
	// req := c.NewRequest(context.Background()).Path("/coffee/hot")

	// var hotCoffees []Coffee
	// if err := c.GetJSON(context.Background(), req, &hotCoffees); err != nil {
	// 	panic(err)
	// }
	// pp.Println(hotCoffees)

	// var icedCoffees []Coffee
	// if err := c.GetJSON(context.Background(), req, &icedCoffees); err != nil {
	// 	panic(err)
	// }
	// pp.Println(icedCoffees)

	go runServer()

	coffee := Coffee{
		ID:          5,
		Title:       "Cappuccino",
		Description: "Made with espresso",
		Ingredients: []string{"Coffee", "Milk", "Sugar"},
	}

	coffeeBytes, _ := json.Marshal(coffee)

	ctx := context.Background()
	req := c.NewRequest().
		Path("/api").
		Method(http.MethodPost).
		Body(coffeeBytes).
		AddHeader("Content-Type", "application/json").
		AddHeader("Accept", "application/json")

	var response ErrorMessage
	err := c.PostJSON(ctx, req, &response)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}

func runServer() {
	e := echo.New()
	e.POST("/api", Test)
	e.Start(":8080")
}

func Test(c echo.Context) error {
	var coffee Coffee
	if err := c.Bind(&coffee); err != nil {
		return err
	}

	return c.JSON(500, ErrorMessage{
		Message: "Internal server error",
		Code:    1000,
	})
}

type ErrorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
