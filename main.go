package main

import "context"

type MockClient struct {
	Client
}

type Mock struct{}

func main() {
	cli := New()

	var mock Mock
	req := cli.NewRequest(context.Background()).
		Path("/orders/%d", 12323)

	if err := cli.GetJSON(context.Background(), req, &mock); err != nil {
		panic(err)
	}
}
