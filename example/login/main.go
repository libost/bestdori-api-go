package main

import (
	"context"

	bestdoriapi "github.com/libost/bestdori-api-go"
)

func main() {
	client, err := bestdoriapi.NewClientWithLogin(context.TODO(), nil, "", "") // Replace with actual username and password
	if err != nil {
		panic(err)
	}
	_ = client
}
