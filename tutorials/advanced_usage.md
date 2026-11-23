# Advanced Usage
In this tutorial, we will explore some advanced usage of the Bestdori API, including fetching recently added songs, events and gachas.

## Fetch What's New
The Bestdori API provides a method to fetch the most recently added songs, events and gachas. This can be useful for keeping your application up to date with the latest content.

### Everything starts with...
You need to get the latest data from the Bestdori API. You can do this by calling the `GetRecent` methods.

```go
recent, err := client.GetRecent(context.TODO())
if err != nil {
    log.Fatal(err)
}
```

### Events

To query the recent event in specific region, you can call the `GetRecentEventsForRegion` method. This will return a list of recent events in the specified region, but it should only contain one event, which is the most recent one.

```go
recentEvents := recent.GetRecentEventsForRegion(api.RegionCN)
if len(recentEvents) == 0 {
    panic("No recent events found")
}
if len(recentEvents) > 1 {
    panic("More than one recent event found") // This should not happen, but we check just in case
}
recentEvent := recentEvents[0]
fmt.Printf("Recent Event Title: %s\n", recentEvent.GetTitle(api.RegionCN)) // example output: "Recent Event Title: Beyond my Imagination"
fmt.Printf("Recent Event Start Time: %s\n", recentEvent.GetStartTime(api.RegionCN).Format("2006-01-02 15:04:05")) // example output: "Recent Event Start Time: 2026-11-23 05:00:00"
fmt.Printf("Recent Event End Time: %s\n", recentEvent.GetEndTime(api.RegionCN).Format("2006-01-02 15:04:05")) // example output: "Recent Event End Time: 2026-11-22 15:59:59"
```

>[!TIP]
>You can always use `GetEvent` method to fetch the details of a recent event.

### Gachas
To query the recent gacha in specific region, you can call the `GetRecentGachasForRegion` method. This will return a list of recent gachas in the specified region.

```go
recentGachas := recent.GetRecentGachasForRegion(api.RegionCN)
if len(recentGachas) == 0 {
    panic("No recent gachas found")
}
recentGacha := recentGachas[0]
fmt.Printf("Recent Gacha Title: %s\n", recentGacha.GetTitle(api.RegionCN)) // example output: "Recent Gacha Title: 高松灯生日纪念 Hopeful birthday！招募"
```

>[!TIP]
>You can always use `GetGacha` method to fetch the details of a recent gacha.

### Songs
To query the recent song in specific region, you can call the `GetRecentSongsForRegion` method. This will return a list of recent songs in the specified region.

```go
recentSongs := recent.GetRecentSongsForRegion(api.RegionCN)
if len(recentSongs) == 0 {
    panic("No recent songs found")
}
recentSong := recentSongs[0]
fmt.Printf("Recent Song Title: %s\n", recentSong.GetTitle(api.RegionCN)) // example output: "Recent Song Title: WHAT AN EXPLOSION"
```

>[!TIP]
>You can always use `GetSong` method to fetch the details of a recent song.

## Summary
In this tutorial, we have learned how to fetch the most recently added songs, events and gachas using the Bestdori API. We covered how to get the latest data, query recent events, gachas and songs in specific regions, and fetch their details. With this knowledge, you can now keep your application up to date with the latest content from the Bestdori API.