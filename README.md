# Go Client 

This library is a go client library and underneath it uses go's `net/http` library. When we send `HTTP` requests via using the `net/http` library. We usually have duplications. We need to create our own retry logics and rate limiting mechanisms. This library is a lightweight helper library to help you use client functionality very clean and good.

## Advantages

- It allows you to write very clean client codes. 
- You don't need to worry about retrying mechanisms, retrying intervals...
- Highly consistent when you have a failure and if your retry count is finally reached the end library provides you a dead-letter mechanism. So whenever you reached retry limit it will use the interface to send this letter to some database, message queue, file etc.
- Provides hands-on rate limiting for you. You only need to worry about the configurations

## Getting Started

Directly use the client struct that library provides you.

```go
package main

import (
	"context"

	"github.com/bilginyuksel/client"
	"github.com/yudai/pp"
)

func main() {
    ctx := context.Background()
    c := client.New(client.WithHost("https://api.sampleapis.com"))
    req := c.NewRequest(ctx).Path("/coffee/hot")

    var hotCoffees []Coffee
    if err := c.GetJSON(ctx, req, &hotCoffees); err != nil {
        panic(err)
    }
    pp.Println(hotCoffees)

    var icedCoffees []Coffee
    if err := c.GetJSON(ctx, req, &icedCoffees); err != nil {
        panic(err)
    }
    pp.Println(icedCoffees)
}

type Coffee struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
}

```

Create a new client struct and embed the library client. Use your new struct with meaningful methods.

```go
package main

import (
	"context"

	"github.com/bilginyuksel/client"
	"github.com/yudai/pp"
)

func main() {
    coffeeClient := NewCoffeeClient(client.WithHost("https://api.sampelapis.com"))
    hotCoffees, err := coffeeClient.GetHotCoffees(context.Background())
    if err != nil {
        panic(err)
    }
    pp.Println(hotCoffees)

    icedCoffees, err := coffeeClient.GetIcedCoffees(context.Background())
    if err != nil {
        panic(err)
    }
    pp.Println(icedCoffees)
}

type Coffee struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
}

// Embed the client that library provides you into the CoffeeClient
type CoffeeClient struct {
	*client.Client
}

func NewCoffeeClient(options client.Option) *MockClient {
	return &CoffeeClient{
		Client: client.New(options),
	}
}

func (c *CoffeeClient) GetHotCoffees(ctx context.Context) ([]Coffee, error) {
    var coffees []Coffee
    req := c.NewRequest(ctx).Path("/coffee/hot")
    return coffees, c.GetJSON(ctx, req, &coffees)
}

func (c *CoffeeClient) GetIcedCoffees(ctx context.Context) ([]Coffee, error) {
    var coffees []Coffee
    req := c.NewRequest(ctx).Path("/coffee/iced")
    return coffees, c.GetJSON(ctx, req, &coffees)
}

```
