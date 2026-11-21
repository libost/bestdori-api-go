# Bestdori API for Go
[![Go Reference](https://pkg.go.dev/badge/github.com/libost/bestdori-api-go.svg)](https://pkg.go.dev/github.com/libost/bestdori-api-go)  
A Go library for interacting with the Bestdori API.

>[!IMPORTANT]
> This is a work in progress and not yet ready for production use.

>[!IMPORTANT]
> This library is not affiliated with the official Bestdori API or its developers. Use at your own risk, as the API may change without notice and this library may not be updated accordingly.

## Usage

- Import the package:

```go
go get github.com/libost/bestdori-api-go
```

## Example

```go
package main

import (
	"context"
	"fmt"
	"time"

	bestdoriapi "github.com/libost/bestdori-api-go"
	"github.com/libost/bestdori-api-go/api"
)

func main() {
	client := bestdoriapi.NewClient("", 10*time.Second)

	card, err := client.GetCard(context.TODO(), "947")
	if err != nil {
		panic(err)
	}
	fmt.Println(card.GetCardName(api.RegionCN))
}
```