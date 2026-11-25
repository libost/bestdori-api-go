# 入门
这是一个简短的教程，介绍如何使用 Bestdori API。

>[!IMPORTANT]
> 在本教程中，我们假设您已经设置好了 Go 环境，并且对 Go 编程有基本的了解。

## 安装
```go
go get github.com/libost/bestdori-api-go
```

## 设置客户端
要与 Bestdori API 交互，您需要创建一个客户端实例。以下是设置方法：
```go
package main

import (    
    bestdoriapi "github.com/libost/bestdori-api-go"
)

func main() {
    client := bestdoriapi.NewClient(nil) // 您可以根据需要传递一个自定义的 http.Client，否则传递 nil 使用默认设置。
}
```

>[!TIP]
> 如果需要，您还可以通过将代理 URL 传递给 `http.Client` 来设置代理。例如：
>```go
>client := bestdoriapi.NewClient(&http.Client{
>    Transport: &http.Transport{
>        Proxy: http.ProxyURL(proxyURL),
>    },
>})
>```

## 总结
在本教程中，我们介绍了如何在 Go 中设置 Bestdori API 客户端。您学习了如何安装该包并创建客户端实例，以开始与 API 交互。

几乎所有您需要获取的内容都封装在 `client` 对象中，因此您可以通过调用其方法开始探索 API。