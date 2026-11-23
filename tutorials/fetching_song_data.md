# Fetching Song Data
In this tutorial, we will learn how to fetch song data using the Bestdori API. The Bestdori API provides a convenient way to access song information for the game BanG Dream! Girls Band Party.

>[!IMPORTANT]
> During this tutorial, we assume that you have already read the [Getting Started](getting_started.md) tutorial and have set up your Go environment.

## Get a Song map
Bestdori API provides a list of all songs in the game. You can fetch this list using the `InitSongs` method.

```go
songsMap, err := client.InitSongs(context.TODO())
if err != nil {
    log.Fatal(err)
}
```

Now you have a map of all songs, here's how you can access a specific song by its ID:

```go
song, err := songsMap.GetSong("306") // Replace "306" with the desired song ID
if err != nil {
    log.Fatal(err)
}
// Always check if the song is available in the region you want to access before fetching its details.
if !song.IsAvailable(api.RegionEN) {
    log.Fatal("Song is not available")
} else {
    // Title of the song in English
    fmt.Println(song.GetSongName(api.RegionEN)) // "SAVIOR OF SONG"
}
If !song.IsAvailable(api.RegionJP) {
    log.Fatal("Song is not available")
} else {
    // Title of the song in Japanese
    fmt.Println(song.GetSongName(api.RegionJP)) // "SAVIOR OF SONG"
}
```

## Get Detailed Song Information
But you may find that the information is not enough, and you want to get more details about the song. You can use the `GetSong` method from the client to fetch detailed information about a specific song.

```go
songDetails, err := client.GetSong(context.TODO(), "306") // Replace "306" with the desired song ID
if err != nil {
    log.Fatal(err)
}
duration := songDetails.GetDuration().Seconds()
fmt.Printf("Song Duration: %.3f seconds\n", duration) // "Song Duration: 90.048 seconds"
```

## Get Something More
Who performed the song? What is the BPM of the song? You can get this information from the `songDetails` object.

```go
bandName, err := songDetails.GetBandName(context.TODO(), api.RegionEN)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Band Name: %s\n", bandName) // "Band Name: RAISE A SUILEN × Yukina"
fmt.Printf("BPM: %d\n", songDetails.GetBPM()) // "BPM: 192"
```

## Metadata
For power players, they are more interested in the metadata of the song, such as efficiency and score factors, as well as skill reliance. You can get this information from the `songDetails` object.

```go
meta, err := songDetails.GetSongMetaFactors(context.TODO(), api.DifficultyExpert)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Efficiency: %d%%\n", math.Round(meta.Efficiency(true, 1.0, 1.1, duration)*100)) // "Efficiency: 435%"
fmt.Printf("Score Factor: %d%%\n", math.Round(meta.ScoreFactor(true, 1.0, 1.1)*100)) // "Score Factor: 869%"
fmt.Printf("Skill Reliance: %d%%\n", math.Round(meta.SkillReliance(true)*100)) // "Skill Reliance: 60%"
```

## Summary
In this tutorial, we learned how to fetch song data using the Bestdori API. We covered how to get a list of songs, access specific song details, and retrieve metadata such as efficiency, score factors, and skill reliance. With this knowledge, you can now explore the Bestdori API further and integrate song data into your applications.