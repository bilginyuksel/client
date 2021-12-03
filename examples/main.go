package main

import (
	"context"
	"log"

	"github.com/bilginyuksel/client"
)

type MockClient struct {
	*client.Client
}

type Mock struct{}

func main() {
	cli := MockClient{Client: client.New()}

	req := cli.NewRequest(context.Background()).
		Path("/orders/%d", 12323)

	var mock Mock
	if err := cli.GetJSON(context.Background(), req, &mock); err != nil {
		panic(err)
	}

	log.Println(mock)
}
