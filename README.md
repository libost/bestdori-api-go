# Bestdori API for Go
[![Go Reference](https://pkg.go.dev/badge/github.com/libost/bestdori-api-go.svg)](https://pkg.go.dev/github.com/libost/bestdori-api-go)
[![Go Version](https://img.shields.io/badge/Go-1.27-blue.svg)](go.mod)
![GitHub Tag](https://img.shields.io/github/v/tag/libost/bestdori-api-go)


A Go library for interacting with the Bestdori API.

[简体中文](README_zh.md)

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

	bestdoriapi "github.com/libost/bestdori-api-go"
	"github.com/libost/bestdori-api-go/api"
)

func main() {
	client := bestdoriapi.NewClient(nil)

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
- More than a API wrapper, it provides additional features such as [event cutoff prediction](tutorials/fetching_event_data.md#Cutoff-Prediction), [gacha emulation](tutorials/fetching_gacha_data.md#Gacha-Emulation), and more.

## Tutorials

Have no idea where to start? Check out the [tutorials](tutorials) directory for step-by-step guides on how to use the Bestdori API.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

This project contains an algorithm ported from [Tsugu BangDream Bot](https://github.com/Yamamoto-2/tsugu-bangdream-bot).
See the [NOTICE](NOTICE) file for details.