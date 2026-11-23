# Fetching Gacha Data
In this tutorial, we will learn how to fetch gacha data using the Bestdori API. The Bestdori API provides a convenient way to access gacha information for the game BanG Dream! Girls Band Party.
>[!IMPORTANT]
> During this tutorial, we assume that you have already read the [Getting Started](getting_started.md) tutorial and have set up your Go environment.

## Get a Gacha map
Bestdori API provides a list of all gachas in the game. You can fetch this list using the `InitGachas` method.

```go
gachasMap, err := client.InitGachas(context.TODO())
if err != nil {
    panic(err)
}
```

Now you have a map of all gachas, here's how you can access a specific gacha by its ID:

```go
gachas, err := gachasMap.GetGacha("1166") // Replace "1166" with the desired gacha ID
if err != nil {
    panic(err)
}
// Always check if the gacha is available in the region you want to access before fetching its details.
if !gachas.IsAvailable(api.RegionJP) {
    panic("Gacha is not available")
}
fmt.Printf("Gacha Name: %s\n", gachas.GetGachaName(api.RegionJP)) // "Gacha Name: 落涙のPouring Rainガチャ"
```

## Get Detailed Gacha Information
But you may find that the information is not enough, and you want to get more details about the gacha. You can use the `GetGacha` method from the client to fetch detailed information about a specific gacha.

```go
gachaDetails, err := client.GetGacha(context.TODO(), "1166") // Replace "1166" with the desired gacha ID
if err != nil {
    panic(err)
}
// Always check if the gacha is available in the region you want to access before fetching its details.
if !gachaDetails.IsAvailable(api.RegionJP) {
    panic("Gacha is not available")
}
fmt.Printf("Gacha Annotation: %s\n", gachaDetails.GetAnnotation(api.RegionJP)) // "Gacha Annotation: ・下記提供割合でメンバーが出現します。\n・出現メンバーはすべて特訓前の状態になります。\n・｢★3以上1人確定｣は｢確定枠提供割合｣をご確認ください。\n・メンバー毎に表示している提供割合は、\n表示上小数点第4位にて四捨五入しているため、\n全てを合算しても100%にならない場合があります。"
```

## Gacha Emulation
For entertainment purposes, you can emulate a gacha pull using the `EmulateGachaDraw` method. This method simulates a gacha pull based on the provided payment method and returns the result. It returns an slice of strings representing the card ID.

```go
paymentMethods := gachaDetails.GetPaymentMethods()
results, err := gachaDetails.EmulateGachaDraw(paymentMethods[0])
if err != nil {
    panic(err)
}
fmt.Printf("Results: %+v\n", results) // example output: "Results: [964 526 883 1448 964 2 911 415 1439 895]"
```

## Summary
In this tutorial, we learned how to fetch gacha data using the Bestdori API. We covered how to fetch a list of gachas, access specific gacha details, and emulate a gacha pull. With this knowledge, you can now explore the gacha data available in the Bestdori API and use it in your applications.