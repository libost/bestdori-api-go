# Bestdori API for Go
[![Go Reference](https://pkg.go.dev/badge/github.com/libost/bestdori-api-go.svg)](https://pkg.go.dev/github.com/libost/bestdori-api-go)
[![Go Version](https://img.shields.io/badge/Go-1.27-blue.svg)](go.mod)
![GitHub Tag](https://img.shields.io/github/v/tag/libost/bestdori-api-go)

一个用于与 Bestdori API 交互的 Go 库。

>[!IMPORTANT]
> 此库仍在开发中，尚未准备好用于生产环境。

>[!IMPORTANT]
> 此库与官方 Bestdori API 或其开发者无关。请自行承担使用风险，因为 API 可能会在没有通知的情况下更改，而此库可能不会相应更新。

## 使用方法
- 导入包：

```go
go get github.com/libost/bestdori-api-go
```

## 示例
>[!NOTE]
> 更多使用示例请参见 [example](example) 目录。

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

## 特性
- 零依赖，仅使用标准库。
- 易于使用且文档齐全。
- 不仅仅是一个 API 包装器，它还提供了额外的功能，例如 [活动分数线预测](tutorials/fetching_event_data.md#Cutoff-Prediction)、[抽卡模拟](tutorials/fetching_gacha_data.md#Gacha-Emulation) 等。

## 教程
不知道如何使用此库？请参阅 [教程](tutorials) 目录中的教程，了解如何获取歌曲、活动、卡牌和抽卡信息，以及如何获取这些内容的详细信息。

## 许可
此项目采用 MIT 许可证 - 有关详细信息，请参阅 [LICENSE](LICENSE) 文件。

## 鸣谢
此项目包含一个移植自 [Tsugu BangDream Bot](https://github.com/Yamamoto-2/tsugu-bangdream-bot) 的算法。
有关详细信息，请参阅 [NOTICE](NOTICE) 。