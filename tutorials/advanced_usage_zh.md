# 高级用法教程
在本教程中，我们将介绍如何使用 Bestdori API 获取最新的歌曲、活动和扭蛋信息。我们将学习如何获取最新数据、查询特定地区的最近活动、扭蛋和歌曲，并获取它们的详细信息。

## 获取最新数据
Bestdori API 提供了一种方法来获取最近添加的歌曲、活动和扭蛋。这对于保持您的应用程序与最新内容同步非常有用。

### 一切都始于...
要获取最新数据，您需要调用 `GetRecent` 方法。

```go
recent, err := client.GetRecent(context.TODO())
if err != nil {
    log.Fatal(err)
}
```

### 活动
要查询特定地区的最近活动，您可以调用 `GetRecentEventsForRegion` 方法。这将返回指定地区的最近活动列表，但它应该只包含一个活动，即最近的一个。

```go
recentEvents := recent.GetRecentEventsForRegion(api.RegionCN)
if len(recentEvents) == 0 {
    panic("No recent events found")
}
if len(recentEvents) > 1 {
    panic("More than one recent event found") // 这不应该发生，但我们还是检查一下以防万一
}
recentEvent := recentEvents[0]
fmt.Printf("Recent Event Title: %s\n", recentEvent.GetTitle(api.RegionCN)) // 示例输出: "Recent Event Title: Beyond my Imagination"
fmt.Printf("Recent Event Start Time: %s\n", recentEvent.GetStartTime(api.RegionCN).Format("2006-01-02 15:04:05")) // 示例输出: "Recent Event Start Time: 2026-11-23 05:00:00"
fmt.Printf("Recent Event End Time: %s\n", recentEvent.GetEndTime(api.RegionCN).Format("2006-01-02 15:04:05")) // 示例输出: "Recent Event End Time: 2026-11-22 15:59:59"
```

>[!TIP]
>您可以始终使用 `GetEvent` 方法来获取最近活动的详细信息。

### 卡池
要查询特定地区的最近卡池，您可以调用 `GetRecentGachasForRegion` 方法。这将返回指定地区的最近卡池列表。

```go
recentGachas := recent.GetRecentGachasForRegion(api.RegionCN)
if len(recentGachas) == 0 {
    panic("No recent gachas found")
}
recentGacha := recentGachas[0]
fmt.Printf("Recent Gacha Title: %s\n", recentGacha.GetTitle(api.RegionCN)) // 示例输出: "Recent Gacha Title: 高松灯生日纪念 Hopeful birthday！招募"
```

>[!TIP]
>您可以始终使用 `GetGacha` 方法来获取最近卡池的详细信息。

### 歌曲
要查询特定地区的最近歌曲，您可以调用 `GetRecentSongsForRegion` 方法。这将返回指定地区的最近歌曲列表。

```go
recentSongs := recent.GetRecentSongsForRegion(api.RegionCN)
if len(recentSongs) == 0 {
    panic("No recent songs found")
}
recentSong := recentSongs[0]
fmt.Printf("Recent Song Title: %s\n", recentSong.GetTitle(api.RegionCN)) // 示例输出: "Recent Song Title: WHAT AN EXPLOSION"
```

>[!TIP]
>您可以始终使用 `GetSong` 方法来获取最近歌曲的详细信息。

## 总结
在本教程中，我们学习了如何使用 Bestdori API 获取最新的歌曲、活动和卡池信息。我们介绍了如何获取最新数据、查询特定地区的最近活动、卡池和歌曲，并获取它们的详细信息。通过这些知识，您现在可以进一步探索 Bestdori API 中可用的最新数据，并将其集成到您的应用程序中。