# Fetching Event Data
In this tutorial, we will learn how to fetch event data using the Bestdori API. The Bestdori API provides a convenient way to access event information for the game BanG Dream! Girls Band Party.
>[!IMPORTANT]
> During this tutorial, we assume that you have already read the [Getting Started](getting_started.md) tutorial and have set up your Go environment.

## Get an Event map
Bestdori API provides a list of all events in the game. You can fetch this list using the `InitEvents` method.

```go
eventsMap, err := client.InitEvents(context.TODO())
if err != nil {
    panic(err)
}
```

Now you have a map of all events, here's how you can access a specific event by its ID:

```go
event, err := eventsMap.GetEvent("311") // Replace "311" with the desired event ID
if err != nil {
    panic(err)
}
if !event.IsAvailable(api.RegionEN) {
    panic("Event is not available")
} 
fmt.Println("Event Name: ", event.GetEventName(api.RegionEN)) // "Event Name: At the Crossroads, Two Stars Smile
if !event.IsAvailable(api.RegionJP) {
    panic("Event is not available")
}
fmt.Println("Event Name: ", event.GetEventName(api.RegionJP)) // "Event Name: 交差点で、ふたつの星が笑う"
```

## Get Detailed Event Information
But you may find that the information is not enough, and you want to get more details about the event. You can use the `GetEvent` method from the client to fetch detailed information about a specific event.

```go
eventDetails, err := client.GetEvent(context.TODO(), "311") // Replace "311" with the desired event ID
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

## Cutoff Prediction
For power players, it is important to know the cutoff for events. The Bestdori API provides a method to predict the cutoff for an event based on the current data.

```go
eventRate, err := eventDetails.GetRates(context.TODO(), api.RegionCN, 50)
if err != nil {
    panic(err)
}
actualCutoff, predictedCutoff, err := eventRate.GetPredictCutoff(context.TODO(), eventDetails)
if err != nil {
    panic(err)
}
// actualCutoff is the actual cutoff for the event
// predictedCutoff is the predicted cutoff for the event
```

## Summary
In this tutorial, we learned how to fetch event data using the Bestdori API. We covered how to get a list of events, access specific event details, and retrieve information about the characters involved in the event. We also learned how to predict the cutoff for an event based on the current data. With this knowledge, you can now explore the Bestdori API further and integrate event data into your applications.