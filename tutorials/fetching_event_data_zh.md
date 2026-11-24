# 获取活动数据
在本教程中，我们将学习如何使用 Bestdori API 获取活动数据。Bestdori API 提供了一种方便的方式来访问 BanG Dream! Girls Band Party 游戏的活动信息。
>[!IMPORTANT]
> 在本教程中，我们假设您已经阅读了 [入门](getting_started.md) 教程，并且已经设置好了 Go 环境。

## 获取活动列表
Bestdori API 提供了游戏中所有活动的列表。您可以使用 `InitEvents` 方法来获取此列表。

```go
eventsMap, err := client.InitEvents(context.TODO())
if err != nil {
    panic(err)
}
```

现在您已经拥有了所有活动的映射，以下是如何通过活动 ID 访问特定活动的方法：

```go
event, err := eventsMap.GetEvent("311") // 将 "311" 替换为所需的活动 ID
if err != nil {
    panic(err)
}
if !event.IsAvailable(api.RegionEN) {
    panic("Event is not available")
} 
// 总是检查活动是否在您想要访问的地区可用，然后再获取其详细信息。
fmt.Println("Event Name: ", event.GetEventName(api.RegionEN)) // "Event Name: At the Crossroads, Two Stars Smile
if !event.IsAvailable(api.RegionJP) {
    panic("Event is not available")
}
fmt.Println("Event Name: ", event.GetEventName(api.RegionJP)) // "Event Name: 交差点で、ふたつの星が笑う"
```

## 获取详细的活动信息
但您可能会发现这些信息还不够，想要获取更多关于活动的详细信息。您可以使用client的 `GetEvent` 方法来获取特定活动的详细信息。

```go
eventDetails, err := client.GetEvent(context.TODO(), "311") // 将 "311" 替换为所需的活动 ID
if err != nil {
    panic(err)
}
eventChars := eventDetails.GetEventCharacters()
if len(eventChars) == 0 {
    panic("No characters found for this event")
}
for _, char := range eventChars {
    character, err := char.GetCharacter(context.TODO())
    if err != nil {
        panic(err)
    }
    fmt.Println("Character Name: ", character.GetCharacterName(api.RegionEN))
} 
/* 
"Character Name: Aya Maruyama" 
"Character Name: Hina Hikawa"
"Character Name: Yukina Minato"
"Character Name: Sayo Hikawa"
"Character Name: Lisa Imai"
*/
```

## 档线预测
对于高水平玩家来说，了解活动的档线是非常重要的。Bestdori API 提供了一种方法，可以根据当前数据预测活动的档线。

```go
eventRate, err := eventDetails.GetRates(context.TODO(), api.RegionCN, 50)
if err != nil {
    panic(err)
}
actualCutoff, predictedCutoff, err := eventRate.GetPredictCutoff(context.TODO(), eventDetails)
if err != nil {
    panic(err)
}
// actualCutoff 是活动的实际档线
// predictedCutoff 是根据当前数据预测的档线
```

## 总结
在本教程中，我们学习了如何使用 Bestdori API 获取活动数据。我们介绍了如何获取活动列表、访问特定活动的详细信息，以及如何预测活动的档线。通过这些知识，您现在可以进一步探索 Bestdori API 中可用的活动数据，并将其集成到您的应用程序中。