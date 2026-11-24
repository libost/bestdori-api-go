# 获取歌曲数据
在本教程中，我们将学习如何使用 Bestdori API 获取歌曲数据。Bestdori API 为我们提供了一个方便的方法来获取歌曲的详细信息，包括歌曲名称、难度、元数据等。

>[!IMPORTANT]
> 在本教程中，我们假设您已阅读了 [入门教程](getting_started_zh.md) 并且已经设置好了 Go 环境。

## 获取歌曲列表
Bestdori API 提供了游戏中所有歌曲的列表。您可以使用 `InitSongs` 方法来获取此列表。

```go
songsMap, err := client.InitSongs(context.TODO())
if err != nil {
    log.Fatal(err)
}
```

现在您已经拥有了所有歌曲的映射，以下是如何通过歌曲 ID 访问特定歌曲的方法：

```go
song, err := songsMap.GetSong("306") // 将 "306" 替换为所需的歌曲 ID
if err != nil {
    log.Fatal(err)
}
// 在获取歌曲详细信息之前，请始终检查该歌曲是否在您想要访问的地区可用。
if !song.IsAvailable(api.RegionEN) {
    log.Fatal("Song is not available")
} else {
    // 获取歌曲的英文名称
    fmt.Println(song.GetSongName(api.RegionEN)) // "SAVIOR OF SONG"
}
If !song.IsAvailable(api.RegionJP) {
    log.Fatal("Song is not available")
} else {
    // 获取歌曲的日文名称
    fmt.Println(song.GetSongName(api.RegionJP)) // "SAVIOR OF SONG"
}
```

## 获取详细的歌曲信息
但您可能会发现这些信息还不够，想要获取更多关于歌曲的详细信息。您可以使用客户端的 `GetSong` 方法来获取特定歌曲的详细信息。

```go
songDetails, err := client.GetSong(context.TODO(), "306") // 将 "306" 替换为所需的歌曲 ID
if err != nil {
    log.Fatal(err)
}
duration := songDetails.GetDuration().Seconds()
fmt.Printf("Song Duration: %.3f seconds\n", duration) // "Song Duration: 90.048 seconds"
```

## 获取更多信息
谁演唱了这首歌？这首歌的 BPM 是多少？您可以从 `songDetails` 对象中获取这些信息。

```go
bandName, err := songDetails.GetBandName(context.TODO(), api.RegionEN)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Band Name: %s\n", bandName) // "Band Name: RAISE A SUILEN × Yukina"
fmt.Printf("BPM: %d\n", songDetails.GetBPM()) // "BPM: 192"
```

## 获取元数据
对于高级玩家来说，他们更关心歌曲的元数据，例如效率和分数，以及技能依赖性。您可以从 `songDetails` 对象中获取这些信息。

```go
meta, err := songDetails.GetSongMetaFactors(context.TODO(), api.DifficultyExpert)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Efficiency: %d%%\n", math.Round(meta.Efficiency(true, 1.0, 1.1, duration)*100)) // "Efficiency: 435%"
fmt.Printf("Score Factor: %d%%\n", math.Round(meta.ScoreFactor(true, 1.0, 1.1)*100)) // "Score Factor: 869%"
fmt.Printf("Skill Reliance: %d%%\n", math.Round(meta.SkillReliance(true)*100)) // "Skill Reliance: 60%"
```

## 总结
在本教程中，我们学习了如何使用 Bestdori API 获取歌曲数据。我们介绍了如何获取歌曲列表、访问特定歌曲的详细信息，以及如何获取元数据，如效率、分数因子和技能依赖性。通过这些知识，您现在可以进一步探索 Bestdori API，并将歌曲数据集成到您的应用程序中。