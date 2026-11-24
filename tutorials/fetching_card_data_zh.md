# 获取卡面数据
在本教程中，我们将学习如何使用 Bestdori API 获取卡面数据。Bestdori API 提供了一种方便的方式来访问 BanG Dream! Girls Band Party 游戏的卡面信息。

>[!IMPORTANT]
> 在本教程中，我们假设您已经阅读了 [入门](getting_started.md) 教程，并且已经设置好了 Go 环境。

## 获取卡面列表
Bestdori API 提供了游戏中所有卡面的列表。您可以使用 `InitCards` 方法来获取此列表。

```go
cardsMap, err := client.InitCards(context.TODO())
if err != nil {
    log.Fatal(err)
}
```

现在您已经拥有了所有卡面的映射，以下是如何通过卡面 ID 访问特定卡面的方法：

```go
card, err := cardsMap.GetCard("947") // 将 "947" 替换为所需的卡面 ID
if err != nil {
    log.Fatal(err)
}
// 总是检查卡面是否在您想要访问的地区可用，然后再获取其详细信息。
if !card.IsAvailable(api.RegionEN) {
    log.Fatal("Card is not available")
} else {
    // 卡面的英文标题
    fmt.Println(card.GetCardName(api.RegionEN)) // Hina's Big Score
}
if !card.IsAvailable(api.RegionJP) {
    log.Fatal("Card is not available")
} else {
    // 卡面的日文标题
    fmt.Println(card.GetCardName(api.RegionJP)) // 日菜的大当たり
}
```

## 获取详细的卡面信息
但您可能会发现这些信息还不够，想要获取更多关于卡面的详细信息。您可以使用客户端的 `GetCard` 方法来获取特定卡面的详细信息。

```go
cardDetails, err := client.GetCard(context.TODO(), "947") // 将 "947" 替换为所需的卡面 ID
if err != nil {
    log.Fatal(err)
}
// 卡面的技能名称（英文）
fmt.Println(cardDetails.GetSkillName(api.RegionEN)) // Boppin' Request♪
// 卡面的技能名称（日文）
fmt.Println(cardDetails.GetSkillName(api.RegionJP)) // るんっ♪とリクエスト
```

>[!TIP]
> 您也可以直接从 `Cards` 结构体中获取特定卡面的详细信息，只需调用该卡面实例上的 `GetCard` 方法即可。
>```go
>cardDetails, err := card.GetCard(context.TODO())
>if err != nil {
>    log.Fatal(err)
>}
>```

## 获取更多内容
如何知道卡面属于谁？您可以使用 `GetCharacter` 方法来获取卡面的角色信息。

```go
character, err := card.GetCharacter(context.TODO())
if err != nil {
    log.Fatal(err)
}
// 卡面的角色名称（英文）
fmt.Println(character.GetCharacterName(api.RegionEN)) // Hina Hikawa
// 卡面的角色名称（日文）
fmt.Println(character.GetCharacterName(api.RegionJP)) // 氷川 日菜
```

## 总结
在本教程中，我们学习了如何使用 Bestdori API 获取卡面数据。我们介绍了如何获取卡面列表、访问特定卡面的详细信息，以及如何获取与卡面相关的角色信息。通过这些知识，您现在可以进一步探索 Bestdori API 中可用的卡面数据，并将其集成到您的应用程序中。