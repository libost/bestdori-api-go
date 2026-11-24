# 获取卡池数据
在本教程中，我们将学习如何使用 Bestdori API 获取卡池数据。Bestdori API 提供了一种方便的方式来访问 BanG Dream! Girls Band Party 游戏的卡池信息。
>[!IMPORTANT]
> 在本教程中，我们假设您已经阅读了 [入门](getting_started.md) 教程，并且已经设置好了 Go 环境。

## 获取卡池列表
Bestdori API 提供了游戏中所有卡池的列表。您可以使用 `InitGachas` 方法来获取此列表。

```go
gachasMap, err := client.InitGachas(context.TODO())
if err != nil {
    panic(err)
}
```

现在您已经拥有了所有卡池的映射，以下是如何通过卡池 ID 访问特定卡池的方法：

```go
gachas, err := gachasMap.GetGacha("1166") // 将 "1166" 替换为所需的卡池 ID
if err != nil {
    panic(err)
}
// 总是检查卡池是否在您想要访问的地区可用，然后再获取其详细信息。
if !gachas.IsAvailable(api.RegionJP) {
    panic("Gacha is not available")
}
fmt.Printf("Gacha Name: %s\n", gachas.GetGachaName(api.RegionJP)) // "Gacha Name: 落涙のPouring Rainガチャ"
```

## 获取详细的卡池信息
但您可能会发现这些信息还不够，想要获取更多关于卡池的详细信息。您可以使用client的 `GetGacha` 方法来获取特定卡池的详细信息。

```go
gachaDetails, err := client.GetGacha(context.TODO(), "1166") // 将 "1166" 替换为所需的卡池 ID
if err != nil {
    panic(err)
}
// 总是检查卡池是否在您想要访问的地区可用，然后再获取其详细信息。
if !gachaDetails.IsAvailable(api.RegionJP) {
    panic("Gacha is not available")
}
fmt.Printf("Gacha Annotation: %s\n", gachaDetails.GetAnnotation(api.RegionJP)) // "Gacha Annotation: ・下記提供割合でメンバーが出現します。\n・出現メンバーはすべて特訓前の状態になります。\n・｢★3以上1人確定｣は｢確定枠提供割合｣をご確認ください。\n・メンバー毎に表示している提供割合は、\n表示上小数点第4位にて四捨五入しているため、\n全てを合算しても100%にならない場合があります。"
```

## 模拟抽卡
出于娱乐目的，您可以使用 `EmulateGachaDraw` 方法来模拟一次抽卡。此方法根据提供的支付方式模拟一次抽卡，并返回结果。它返回一个字符串切片，表示卡牌 ID。

```go
paymentMethods := gachaDetails.GetPaymentMethods()
results, err := gachaDetails.EmulateGachaDraw(paymentMethods[0])
if err != nil {
    panic(err)
}
fmt.Printf("Results: %+v\n", results) // 示例输出: "Results: [964 526 883 1448 964 2 911 415 1439 895]"
```

## 总结
在本教程中，我们学习了如何使用 Bestdori API 获取卡池数据。我们介绍了如何获取卡池列表、访问特定卡池的详细信息，以及如何模拟一次抽卡。通过这些知识，您现在可以进一步探索 Bestdori API 中可用的卡池数据，并将其集成到您的应用程序中。