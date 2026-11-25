# Getting Started

This is a simple guide to get you started with the Bestdori API.  
>[!IMPORTANT]
>During this tutorial, we assume you have already set up your Go environment and have a basic understanding of Go programming.

## Installation
```go
go get github.com/libost/bestdori-api-go
```

## Setting Up the Client
To interact with the Bestdori API, you need to create a client instance. Here's how to do it:

```go
package main

import (
    "time"
    
    bestdoriapi "github.com/libost/bestdori-api-go"
)

func main() {
    client := bestdoriapi.NewClient(nil) // You can pass a custom http.Client if needed, otherwise pass nil for default settings.
}
```

>[!TIP]
> You can also set a proxy if needed by passing the proxy URL inside a `http.Client` when creating the client. For example:
>```go
>client := bestdoriapi.NewClient(&http.Client{
>    Transport: &http.Transport{
>        Proxy: http.ProxyURL(proxyURL),
>    },
>})
>```

## Summary
In this tutorial, we covered the basics of setting up the Bestdori API client in Go. You learned how to install the package and create a client instance to start interacting with the API. 

Almost everything you would need to get is wrapped in the `client` object, so you can start exploring the API by calling its methods.