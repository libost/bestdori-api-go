# Bestdori API for Go
[![Go Reference](https://pkg.go.dev/badge/github.com/libost/bestdori-api-go.svg)](https://pkg.go.dev/github.com/libost/bestdori-api-go)
[![Go Version](https://img.shields.io/badge/Go-1.27-blue.svg)](go.mod)

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
>[!NOTE]
> See the [example](example) directory for more usage examples.  

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
	if !card.IsAvailable(api.RegionCN) {
		fmt.Println("Card is not available in CN region")
		return
	}
	fmt.Println("Card is available in CN region")
	fmt.Println(card.GetCardName(api.RegionCN))
}
```

## Features

- Zero dependencies, only the standard library is used.
- Easy to use and well-documented.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.