# Fetching Card Data
In this tutorial, we will learn how to fetch card data using the Bestdori API. The Bestdori API provides a convenient way to access card information for the game BanG Dream! Girls Band Party.

>[!IMPORTANT]
> During this tutorial, we assume that you have already read the [Getting Started](getting_started.md) tutorial and have set up your Go environment.

## Get a Card map
Bestdori API provides a list of all cards in the game. You can fetch this list using the `InitCards` method.

```go
cardsMap, err := client.InitCards(context.TODO())
if err != nil {
    log.Fatal(err)
}
```

Now you have a map of all cards, here's how you can access a specific card by its ID:

```go
card, err := cardsMap.GetCard("947") // Replace "947" with the desired card ID
if err != nil {
    log.Fatal(err)
}
// Always check if the card is available in the region you want to access before fetching its details.
if !card.IsAvailable(api.RegionEN) {
    log.Fatal("Card is not available")
} else {
    // Title of the card in English
    fmt.Println(card.GetCardName(api.RegionEN)) // Hina's Big Score
}
if !card.IsAvailable(api.RegionJP) {
    log.Fatal("Card is not available")
} else {
    // Title of the card in Japanese
    fmt.Println(card.GetCardName(api.RegionJP)) // 日菜的大当たり
}
```

## Get Detailed Card Information
But you may find that the information is not enough, and you want to get more details about the card. You can use the `GetCard` method from the client to fetch detailed information about a specific card.

```go
cardDetails, err := client.GetCard(context.TODO(), "947") // Replace "947" with the desired card ID
if err != nil {
    log.Fatal(err)
}
// Skill name of the card in English
fmt.Println(cardDetails.GetSkillName(api.RegionEN)) // Boppin' Request♪
// Skill name of the card in Japanese
fmt.Println(cardDetails.GetSkillName(api.RegionJP)) // るんっ♪とリクエスト
```

>[!TIP]
> You can also fetch detailed information about a card directly from `Cards` struct by calling the `GetCard` method on a specific card instance.
>```go
>cardDetails, err := card.GetCard(context.TODO())
>if err != nil {
>    log.Fatal(err)
>}
>```

## Get Something More
How do I know who the card belongs to? You can use the `GetCharacter` method to fetch the character of the card.

```go
character, err := card.GetCharacter(context.TODO())
if err != nil {
    log.Fatal(err)
}
// Character name of the card in English
fmt.Println(character.GetCharacterName(api.RegionEN)) // Hina Hikawa
// Character name of the card in Japanese
fmt.Println(character.GetCharacterName(api.RegionJP)) // 氷川 日菜
```

## Summary
In this tutorial, we covered how to fetch card data using the Bestdori API. We learned how to get a map of all cards, access specific card details, and fetch additional information such as the character associated with a card. With this knowledge, you can now explore the Bestdori API further and retrieve the card data you need for your applications.